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

type BuildingUsecase struct {
	buildRepo   repo.BuildingRepository
	assetHelper helper.AssetHelper
	userClient  userPb.UserServiceClient
}

func NewBuildingUsecase(r repo.BuildingRepository, ah helper.AssetHelper, uc userPb.UserServiceClient) *BuildingUsecase {
	return &BuildingUsecase{
		buildRepo:   r,
		assetHelper: ah,
		userClient:  uc,
	}
}

func (u *BuildingUsecase) CreateBuilding(ctx context.Context, req *pb.CreateBuildingRequest) (*pb.BuildingResponse, error) {
	userID, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	building := mapper.ToBuildingDomain(req, userID)

	if err := u.buildRepo.CreateBuilding(ctx, building); err != nil {
		return nil, err
	}

	return &pb.BuildingResponse{
		Building: mapper.ToBuildingProto(building),
	}, nil
}

func (u *BuildingUsecase) GetBuilding(ctx context.Context, req *pb.GetAssetRequest) (*pb.BuildingArrayResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	buildings, err := u.buildRepo.GetBuilding(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &pb.BuildingArrayResponse{
		Success:  true,
		Building: mapper.ToBuildingProtoSlice(buildings),
	}, nil
}

func (u *BuildingUsecase) GetBuildingByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.BuildingArrayResponse, error) {
	ids := utils.ParseUUIDs(req.Ids)
	bu, err := u.buildRepo.GetBuildingByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &pb.BuildingArrayResponse{
		Building: mapper.ToBuildingProtoSlice(bu),
	}, nil
}

func (u *BuildingUsecase) GetBatchBuildingByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.BuildingArrayResponse, error) {
	ids := utils.ParseUUIDs(req.Ids)
	bu, err := u.buildRepo.GetBatchBuildingByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &pb.BuildingArrayResponse{
		Building: mapper.ToBuildingProtoSlice(bu),
	}, nil
}

func (u *BuildingUsecase) GetBuildingByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.BuildingResponse, error) {
	id, err := utils.ParseID(req.Id)
	if err != nil {
		return nil, err
	}

	building, err := u.buildRepo.GetBuildingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.BuildingResponse{
		Success:  true,
		Building: mapper.ToBuildingProto(building),
	}, nil
}

func (u *BuildingUsecase) UpdateBuilding(ctx context.Context, req *pb.UpdateBuildingRequest) (*pb.BuildingResponse, error) {
	if req.Building == nil {
		return nil, errors.New("building data is required")
	}

	id, uid, err := utils.ValidateIDs(req.Id, req.Building.UserId)
	if err != nil {
		return nil, err
	}

	building, err := u.buildRepo.GetBuildingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := helper.ApplyUpdateBuildingFields(req, building); err != nil {
		return nil, err
	}

	syncParams := domain.FileSyncParams{
		UserID:        uid,
		EntityID:      id,
		EntityType:    "building",
		NewFiles:      req.NewFiles,
		DeleteFileIDs: req.DeleteFileIds,
	}

	if err := u.assetHelper.SyncFiles(ctx, syncParams); err != nil {
		return nil, err
	}

	updatedBuilding, err := u.buildRepo.UpdateBuilding(
		ctx, building, utils.ParseUUIDs(req.LandIds), utils.ParseUUIDs(req.DeleteLandIds), utils.ParseUUIDs(req.InsIds), utils.ParseUUIDs(req.DeleteInsIds),
	)
	if err != nil {
		return nil, err
	}

	return &pb.BuildingResponse{
		Success:  true,
		Building: mapper.ToBuildingProto(updatedBuilding),
	}, nil
}

func (u *BuildingUsecase) DeleteBuilding(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	if _, err = u.buildRepo.GetBuildingByID(ctx, id); err != nil {
		return nil, err
	}

	if err := u.buildRepo.SoftDeleteBuilding(ctx, id, uid); err != nil {
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

func (u *BuildingUsecase) CleanupExpiredBuildings(ctx context.Context) error {
	cutoffTime := time.Now().AddDate(0, 0, -7)
	expiredBuilding, err := u.buildRepo.GetExpiredBuilding(ctx, cutoffTime)
	if err != nil {
		return err
	}

	if len(expiredBuilding) == 0 {
		return nil
	}

	for _, bu := range expiredBuilding {
		u.assetHelper.CleanupResource(ctx, bu.ID, bu.Files, func(id uuid.UUID) error {
			return u.buildRepo.HardDeleteBuilding(ctx, id)
		})
	}

	return nil
}
