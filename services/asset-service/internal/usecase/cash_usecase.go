package usecase

import (
	"context"
	"errors"
	"log"
	"time"
	"wealth-vault/asset-service/internal/domain"
	repo "wealth-vault/asset-service/internal/repository/interface"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	userPb "wealth-vault/asset-service/pkg/pb/proto/user"
	"wealth-vault/asset-service/pkg/utils"
	helper "wealth-vault/asset-service/pkg/utils/helper"
	"wealth-vault/asset-service/pkg/utils/mapper"

	"github.com/google/uuid"
)

type CashUsecase struct {
	cashRepo    repo.CashRepository
	assetHelper helper.AssetHelper
	userClient  userPb.UserServiceClient
}

func NewCashUsecase(r repo.CashRepository, ah helper.AssetHelper, uc userPb.UserServiceClient) *CashUsecase {
	return &CashUsecase{
		cashRepo:    r,
		assetHelper: ah,
		userClient:  uc,
	}
}

func (u *CashUsecase) CreateCash(ctx context.Context, req *pb.CreateCashRequest) (*pb.CashResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	cash := mapper.ToCashDomain(req, uid)
	if err := u.cashRepo.CreateCash(ctx, cash); err != nil {
		return nil, err
	}

	return &pb.CashResponse{
		Cash: mapper.ToCashProto(cash),
	}, nil
}

func (u *CashUsecase) GetCash(ctx context.Context, req *pb.GetAssetRequest) (*pb.CashArrayResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	cash, err := u.cashRepo.GetCash(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &pb.CashArrayResponse{
		Success: true,
		Cash:    mapper.ToCashProtoSlice(cash),
	}, nil
}

func (u *CashUsecase) GetCashByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.CashArrayResponse, error) {
	ids := utils.ParseUUIDs(req.Ids)

	cash, err := u.cashRepo.GetCashByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &pb.CashArrayResponse{
		Cash: mapper.ToCashProtoSlice(cash),
	}, nil
}

func (u *CashUsecase) GetBatchCashByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.CashArrayResponse, error) {
	ids := utils.ParseUUIDs(req.Ids)

	cash, err := u.cashRepo.GetBatchCashByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &pb.CashArrayResponse{
		Cash: mapper.ToCashProtoSlice(cash),
	}, nil
}

func (u *CashUsecase) GetCashByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.CashResponse, error) {
	id, err := utils.ParseID(req.Id)
	if err != nil {
		return nil, err
	}

	cash, err := u.cashRepo.GetCashByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.CashResponse{
		Success: true,
		Cash:    mapper.ToCashProto(cash),
	}, nil
}

func (u *CashUsecase) UpdateCash(ctx context.Context, req *pb.UpdateCashRequest) (*pb.CashResponse, error) {
	if req.Cash == nil {
		return nil, errors.New("cash data is required")
	}

	id, uid, err := utils.ValidateIDs(req.Id, req.Cash.UserId)
	if err != nil {
		return nil, err
	}

	cash, err := u.cashRepo.GetCashByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := helper.ApplyUpdateCashFields(req, cash); err != nil {
		return nil, err
	}

	syncParams := domain.FileSyncParams{
		UserID:        uid,
		EntityID:      id,
		EntityType:    "cash",
		NewFiles:      req.NewFiles,
		DeleteFileIDs: req.DeleteFileIds,
	}

	if err := u.assetHelper.SyncFiles(ctx, syncParams); err != nil {
		return nil, err
	}

	updatedCash, err := u.cashRepo.UpdateCash(ctx, cash)
	if err != nil {
		return nil, err
	}

	return &pb.CashResponse{
		Success: true,
		Cash:    mapper.ToCashProto(updatedCash),
	}, nil
}

func (u *CashUsecase) DeleteCash(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	if _, err = u.cashRepo.GetCashByID(ctx, id); err != nil {
		return nil, err
	}

	if err := u.cashRepo.SoftDeleteCash(ctx, id, uid); err != nil {
		return nil, err
	}

	go func() {
		bgCtx := context.Background()

		_, err := u.userClient.MarkAssetMessagesDeleted(bgCtx, &userPb.MarkAssetDeletedRequest{
			AssetId: id.String(),
		})

		if err != nil {
			log.Printf("⚠️ Failed to notify User Service via gRPC: %v", err)
		}
	}()

	return &pb.DeleteAssetResponse{
		Success: true,
	}, nil
}

func (u *CashUsecase) CleanupExpiredCashes(ctx context.Context) error {
	cutoffTime := time.Now().AddDate(0, 0, -7)
	expiredCash, err := u.cashRepo.GetExpiredCash(ctx, cutoffTime)
	if err != nil {
		return err
	}

	if len(expiredCash) == 0 {
		return nil
	}

	for _, c := range expiredCash {
		u.assetHelper.CleanupResource(ctx, c.ID, c.Files, func(id uuid.UUID) error {
			return u.cashRepo.HardDeleteCash(ctx, id)
		})
	}

	return nil
}
