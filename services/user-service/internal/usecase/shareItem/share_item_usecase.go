package usecase

import (
	"context"
	"wealth-vault/user-service/internal/event"
	repo "wealth-vault/user-service/internal/repository/interface"
	"wealth-vault/user-service/pkg/mail"
	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ShareItemUsecase struct {
	itemRepo    repo.ShareItemRepository
	groupRepo   repo.GroupRepository
	userRepo    repo.UserRepository
	msgRepo     repo.MsgRepository
	assetClient assetPb.AssetServiceClient
	mailclient  mail.NotificationClient
	publisher   event.EventPublisher
}

func NewShareItemUsecase(
	r repo.ShareItemRepository,
	g repo.GroupRepository,
	u repo.UserRepository,
	m repo.MsgRepository,
	assetClient assetPb.AssetServiceClient,
	mail mail.NotificationClient,
	e event.EventPublisher,
) *ShareItemUsecase {
	return &ShareItemUsecase{
		itemRepo:    r,
		groupRepo:   g,
		userRepo:    u,
		assetClient: assetClient,
		mailclient:  mail,
		publisher:   e,
		msgRepo:     m,
	}
}

const (
	AssetTypeBuilding   = "building"
	AssetTypeLand       = "land"
	AssetTypeAccount    = "account"
	AssetTypeCash       = "cash"
	AssetTypeInsurance  = "insurance"
	AssetTypeInvestment = "investment"
	AssetTypeLiability  = "liability"

	LogTypeActivity = "ACTIVITY"
	LogTypeSystem   = "SYSTEM"

	MsgTypeAssetCard   = "ASSET_CARD"
	MsgTypeSystemAlert = "SYSTEM_ALERT"

	TopicItemShared       = "noti.item.shared"
	TopicGroupActivity    = "noti.group.activity"
	TopicGroupMemberAdded = "noti.group.member.added"
	TopicAccessGranted    = "noti.access.granted"
)

func (u *ShareItemUsecase) DeleteAllReferencesByEntityID(ctx context.Context, req *pb.DeleteByEntityRequest) (*pb.DeleteByEntityResponse, error) {
	id, err := utils.ParseUUID(req.EntityId)
	if err != nil {
		return nil, err
	}

	if err := u.itemRepo.DeleteAllReferencesByEntityID(ctx, id); err != nil {
		return nil, err
	}

	return &pb.DeleteByEntityResponse{
		Success: true,
	}, nil
}

func (u *ShareItemUsecase) GetItemSharedTargets(ctx context.Context, req *pb.GetItemSharedTargetsRequest) (*pb.GetItemSharedTargetsResponse, error) {
	itemId, err := utils.ParseUUID(req.ItemId)
	if err != nil {
		return nil, err
	}

	userId, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	targets, err := u.itemRepo.GetItemSharedTargets(ctx, userId, itemId, req.ItemType)
	if err != nil {
		return nil, err
	}

	resp := &pb.GetItemSharedTargetsResponse{}
	for _, g := range targets.Groups {
		resp.Groups = append(resp.Groups, &pb.SharedGroupInfo{
			GroupId:     g.GroupID,
			GroupName:   g.GroupName,
			GroupImage:  g.GroupImage,
			MemberCount: g.MemberCount,
			SharedAt:    timestamppb.New(g.SharedAt),
		})
	}

	for _, f := range targets.Friends {
		resp.Friends = append(resp.Friends, &pb.SharedFriendInfo{
			FriendId:     f.FriendID,
			Username:     f.Username,
			ProfileImage: f.ProfileImage,
			SharedAt:     timestamppb.New(f.SharedAt),
		})
	}

	for _, e := range targets.Emails {
		resp.Emails = append(resp.Emails, &pb.SharedEmailInfo{
			Email:    e.Email,
			SharedAt: timestamppb.New(e.SharedAt),
			IsSent:   e.IsSent,
		})
	}

	return resp, nil
}

func (u *ShareItemUsecase) GetSharedItemIDs(ctx context.Context, req *pb.GetSharedItemIDsRequest) (*pb.GetSharedItemIDsResponse, error) {
	userID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	targetID, err := utils.ParseUUID(req.TargetId)
	if err != nil {
		return nil, err
	}

	ids, err := u.itemRepo.GetSharedItemIDs(ctx, userID, targetID, req.TargetType)
	if err != nil {
		return nil, err
	}

	return &pb.GetSharedItemIDsResponse{
		ItemIds: ids,
	}, nil
}

func (u *ShareItemUsecase) GetAllSharedItemIDsByUser(ctx context.Context, req *pb.GetAllSharedItemIDsByUserRequest) (*pb.GetAllSharedItemIDsByUserResponse, error) {
	userID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	ids, err := u.itemRepo.GetAllSharedItemIDsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &pb.GetAllSharedItemIDsByUserResponse{
		ItemIds: ids,
	}, nil
}
