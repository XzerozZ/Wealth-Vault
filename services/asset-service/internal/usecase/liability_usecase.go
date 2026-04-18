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

type LiabilityUsecase struct {
	liaRepo     repo.LiabilityRepository
	assetHelper helper.AssetHelper
	userClient  userPb.UserServiceClient
}

func NewLiabilityUsecase(r repo.LiabilityRepository, ah helper.AssetHelper, uc userPb.UserServiceClient) *LiabilityUsecase {
	return &LiabilityUsecase{
		liaRepo:     r,
		assetHelper: ah,
		userClient:  uc,
	}
}

func (u *LiabilityUsecase) CreateLiability(ctx context.Context, req *pb.CreateLiabilityRequest) (*pb.LiabilityResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	liability := mapper.ToLiabilityDomain(req, uid)

	if err := u.liaRepo.CreateLiability(ctx, liability); err != nil {
		return nil, err
	}

	return &pb.LiabilityResponse{
		Success:   true,
		Liability: mapper.ToLiabilityProto(liability),
	}, nil
}

func (u *LiabilityUsecase) GetLiability(ctx context.Context, req *pb.GetLiabilityRequest) (*pb.LiabilityArrayResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	lias, err := u.liaRepo.GetLiability(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &pb.LiabilityArrayResponse{
		Success:   true,
		Liability: mapper.ToLiabilityProtoSlice(lias),
	}, nil
}

func (u *LiabilityUsecase) GetLiabilityByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.LiabilityArrayResponse, error) {
	ids := utils.ParseUUIDs(req.Ids)

	lias, err := u.liaRepo.GetLiabilityByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &pb.LiabilityArrayResponse{
		Liability: mapper.ToLiabilityProtoSlice(lias),
	}, nil
}

func (u *LiabilityUsecase) GetBatchLiabilityByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.LiabilityArrayResponse, error) {
	ids := utils.ParseUUIDs(req.Ids)

	lias, err := u.liaRepo.GetBatchLiabilityByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &pb.LiabilityArrayResponse{
		Liability: mapper.ToLiabilityProtoSlice(lias),
	}, nil
}

func (u *LiabilityUsecase) GetLiabilityByID(ctx context.Context, req *pb.GetLiabilityByIDRequest) (*pb.LiabilityResponse, error) {
	id, err := utils.ParseID(req.Id)
	if err != nil {
		return nil, err
	}

	lia, err := u.liaRepo.GetLiabilityByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.LiabilityResponse{
		Success:   true,
		Liability: mapper.ToLiabilityProto(lia),
	}, nil
}

func (u *LiabilityUsecase) UpdateLiability(ctx context.Context, req *pb.UpdateLiabilityRequest) (*pb.LiabilityResponse, error) {
	if req.Liability == nil {
		return nil, errors.New("liability data is required")
	}

	id, uid, err := utils.ValidateIDs(req.Id, req.Liability.CreatedBy)
	if err != nil {
		return nil, err
	}

	lia, err := u.liaRepo.GetLiabilityByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := helper.ApplyUpdateFields(req, lia); err != nil {
		return nil, err
	}

	syncParams := domain.FileSyncParams{
		UserID:        uid,
		EntityID:      id,
		EntityType:    "liability",
		NewFiles:      req.NewFiles,
		DeleteFileIDs: req.DeleteFileIds,
	}

	if err := u.assetHelper.SyncFiles(ctx, syncParams); err != nil {
		return nil, err
	}

	updatedLia, err := u.liaRepo.UpdateLiability(ctx, lia)
	if err != nil {
		return nil, err
	}

	return &pb.LiabilityResponse{
		Success:   true,
		Liability: mapper.ToLiabilityProto(updatedLia),
	}, nil
}

func (u *LiabilityUsecase) DeleteLiability(ctx context.Context, req *pb.DeleteLiabilityRequest) (*pb.DeleteLiabilityResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	if _, err = u.liaRepo.GetLiabilityByID(ctx, id); err != nil {
		return nil, err
	}

	if err := u.liaRepo.SoftDeleteLiability(ctx, id, uid); err != nil {
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

	return &pb.DeleteLiabilityResponse{
		Success: true,
	}, nil
}

func (u *LiabilityUsecase) CleanupExpiredLiabilities(ctx context.Context) error {
	cutoffTime := time.Now().AddDate(0, 0, -7)
	expiredLia, err := u.liaRepo.GetExpiredLiability(ctx, cutoffTime)
	if err != nil {
		return err
	}

	if len(expiredLia) == 0 {
		return err
	}

	for _, l := range expiredLia {
		u.assetHelper.CleanupResource(ctx, l.ID, l.Files, func(id uuid.UUID) error {
			return u.liaRepo.HardDeleteLiability(ctx, id)
		})
	}

	return nil
}
