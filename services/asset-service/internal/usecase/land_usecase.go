package usecase

import (
	"context"
	"errors"
	"time"
	"wealth-vault/asset-service/internal/domain"
	repo "wealth-vault/asset-service/internal/repository/interface"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils"
	helper "wealth-vault/asset-service/pkg/utils/helper"
	"wealth-vault/asset-service/pkg/utils/mapper"

	"github.com/google/uuid"
)

type LandUsecase struct {
	landRepo    repo.LandRepository
	assetHelper helper.AssetHelper
}

func NewLandUsecase(r repo.LandRepository, ah helper.AssetHelper) LandUsecase {
	return LandUsecase{
		landRepo:    r,
		assetHelper: ah,
	}
}

func (u *LandUsecase) CreateLand(ctx context.Context, req *pb.CreateLandRequest) (*pb.LandResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	land := mapper.ToLandDomain(req, uid)

	if err := u.landRepo.CreateLand(ctx, land); err != nil {
		return nil, err
	}

	return &pb.LandResponse{
		Land: mapper.ToLandProto(land),
	}, nil
}

func (u *LandUsecase) GetLand(ctx context.Context, req *pb.GetAssetRequest) (*pb.LandArrayResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	lands, err := u.landRepo.GetLand(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &pb.LandArrayResponse{
		Success: true,
		Land:    mapper.ToLandProtoSlice(lands),
	}, nil
}

func (u *LandUsecase) GetLandByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.LandArrayResponse, error) {
	ids := utils.ParseUUIDs(req.Ids)

	lands, err := u.landRepo.GetLandByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &pb.LandArrayResponse{
		Land: mapper.ToLandProtoSlice(lands),
	}, nil
}

func (u *LandUsecase) GetBatchLandByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.LandArrayResponse, error) {
	ids := utils.ParseUUIDs(req.Ids)

	lands, err := u.landRepo.GetBatchLandByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &pb.LandArrayResponse{
		Land: mapper.ToLandProtoSlice(lands),
	}, nil
}

func (u *LandUsecase) GetLandByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.LandResponse, error) {
	id, err := utils.ParseID(req.Id)
	if err != nil {
		return nil, err
	}

	land, err := u.landRepo.GetLandByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.LandResponse{
		Success: true,
		Land:    mapper.ToLandProto(land),
	}, nil
}

func (u *LandUsecase) UpdateLand(ctx context.Context, req *pb.UpdateLandRequest) (*pb.LandResponse, error) {
	if req.Land == nil {
		return nil, errors.New("land data is required")
	}

	id, uid, err := utils.ValidateIDs(req.Id, req.Land.UserId)
	if err != nil {
		return nil, err
	}

	land, err := u.landRepo.GetLandByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := helper.ApplyUpdateLandFields(req, land); err != nil {
		return nil, err
	}

	syncParams := domain.FileSyncParams{
		UserID:        uid,
		EntityID:      id,
		EntityType:    "land",
		NewFiles:      req.NewFiles,
		DeleteFileIDs: req.DeleteFileIds,
	}

	if err := u.assetHelper.SyncFiles(ctx, syncParams); err != nil {
		return nil, err
	}

	updatedLand, err := u.landRepo.UpdateLand(
		ctx, land, utils.ParseUUIDs(req.BuildingIds), utils.ParseUUIDs(req.DeleteBuildingIds),
	)
	if err != nil {
		return nil, err
	}

	return &pb.LandResponse{
		Success: true,
		Land:    mapper.ToLandProto(updatedLand),
	}, nil
}

func (u *LandUsecase) DeleteLand(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	if _, err = u.landRepo.GetLandByID(ctx, id); err != nil {
		return nil, err
	}

	if err := u.landRepo.SoftDeleteLand(ctx, id, uid); err != nil {
		return nil, err
	}

	return &pb.DeleteAssetResponse{
		Success: true,
	}, nil
}

func (u *LandUsecase) CleanupExpiredLands(ctx context.Context) error {
	cutoffTime := time.Now().AddDate(0, 0, -7)
	expiredLand, err := u.landRepo.GetExpiredLand(ctx, cutoffTime)
	if err != nil {
		return err
	}

	if len(expiredLand) == 0 {
		return nil
	}

	for _, l := range expiredLand {
		u.assetHelper.CleanupResource(ctx, l.ID, l.Files, func(id uuid.UUID) error {
			return u.landRepo.HardDeleteLand(ctx, id)
		})
	}

	return nil
}
