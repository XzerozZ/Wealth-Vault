package usecase

import (
	"context"
	"fmt"
	"log"
	"time"
	"wealth-vault/user-service/internal/domain"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (u *ShareItemUsecase) GetSharedIteminFriend(ctx context.Context, req *pb.GetFriendItemRequest) (*pb.GetFriendItemsResponse, error) {
	uid, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	fid, err := utils.ParseUUID(req.FriendId)
	if err != nil {
		return nil, err
	}

	items, err := u.itemRepo.GetSharedIteminFriend(ctx, fid, uid)
	if err != nil {
		return nil, err
	}

	var summaries []domain.SharedItemSummary
	for _, item := range items {
		summaries = append(summaries, domain.SharedItemSummary{EntityID: item.EntityID.String(), EntityType: item.EntityType})
	}

	previewMap, err := u.FetchAssetPreviews(ctx, summaries)
	if err != nil {
		return nil, err
	}

	var responseItems []*pb.FriendItemDetail
	for _, item := range items {
		responseItems = append(responseItems, &pb.FriendItemDetail{
			FriendItemId: item.ID.String(),
			SharedBy:     item.OwnerID.String(),
			SharedAt:     timestamppb.New(item.CreatedAt),
			Type:         item.EntityType,
			AssetDetail:  previewMap[item.EntityID.String()],
		})
	}
	return &pb.GetFriendItemsResponse{Items: responseItems}, nil
}

func (u *ShareItemUsecase) UnsharedIteminFriend(ctx context.Context, req *pb.UnshareItemRequest) (*pb.ShareItemResponse, error) {
	itemID, err := utils.ParseUUID(req.ItemId)
	if err != nil {
		return nil, err
	}

	userID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	if err := u.itemRepo.DeleteIteminFriend(ctx, itemID, userID); err != nil {
		return nil, err
	}

	go func() {
		if err := u.msgRepo.MarkAssetMessageAsDeleted(context.Background(), itemID); err != nil {
			log.Printf("Failed to mark asset message as deleted: %v", err)
		}
	}()

	return &pb.ShareItemResponse{Finish: true}, nil
}

func (u *ShareItemUsecase) BatchShareAssets(ctx context.Context, req domain.BatchShareRequest) error {
	existingMap, err := u.itemRepo.GetExistingSharedMap(ctx, req.OwnerID, req.TargetID)
	if err != nil {
		return fmt.Errorf("failed to fetch existing shares: %v", err)
	}

	var newItems []domain.FriendItem
	now := time.Now()

	processIDs := func(ids []string, entityType string) {
		for _, idStr := range ids {
			if idStr == "" {
				continue
			}
			if !existingMap[fmt.Sprintf("%s:%s", entityType, idStr)] {
				if uid, err := uuid.Parse(idStr); err == nil {
					newItems = append(newItems, domain.FriendItem{OwnerID: req.OwnerID, FriendID: req.TargetID, EntityType: entityType, EntityID: uid, ShareAt: now})
				}
			}
		}
	}

	processIDs(req.AccountIDs, AssetTypeAccount)
	processIDs(req.BuildingIDs, AssetTypeBuilding)
	processIDs(req.CashIDs, AssetTypeCash)
	processIDs(req.InsuranceIDs, AssetTypeInsurance)
	processIDs(req.InvestmentIDs, AssetTypeInvestment)
	processIDs(req.LandIDs, AssetTypeLand)
	processIDs(req.LiabilityIDs, AssetTypeLiability)

	if len(newItems) > 0 {
		if err := u.itemRepo.ShareItemtoFriend(ctx, newItems); err != nil {
			return fmt.Errorf("failed to batch insert items: %v", err)
		}

		go func() {
			bgCtx := context.Background()
			senderName := "Unknown"
			if userProfile, err := u.userRepo.GetUser(bgCtx, req.OwnerID); err == nil {
				senderName = userProfile.Username
			}

			evt := domain.ItemSharedEvent{
				AssetID:       "",
				SenderID:      req.OwnerID.String(),
				SenderName:    senderName,
				TargetUserIDs: []string{req.TargetID.String()},
				OccurredAt:    time.Now().Unix(),
			}

			u.publisher.Publish("noti.item.shared", evt)
		}()
	}
	return nil
}

func (u *ShareItemUsecase) GetItemsSharedByFriend(ctx context.Context, req *pb.GetItemsSharedByFriendRequest) (*pb.GetItemsSharedByFriendResponse, error) {
	uid, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	friendID, err := utils.ParseUUID(req.FriendId)
	if err != nil {
		return nil, err
	}

	items, err := u.itemRepo.GetItemsSharedByFriend(ctx, uid, friendID)
	if err != nil {
		return nil, err
	}

	previewMap, err := u.FetchAssetPreviews(ctx, items)
	if err != nil {
		return nil, err
	}

	var responseItems []*pb.SharedAssetPreview
	for _, item := range items {
		if preview, ok := previewMap[item.EntityID]; ok {
			responseItems = append(responseItems, &pb.SharedAssetPreview{Id: item.EntityID, Type: item.EntityType, AssetDetail: preview})
		}
	}
	return &pb.GetItemsSharedByFriendResponse{AssetDetail: responseItems}, nil
}
