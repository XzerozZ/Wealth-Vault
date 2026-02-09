package usecase

import (
	"context"
	"errors"
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
	buildRepo  repo.BuildingRepository
	fileRepo   repo.FileRepository
	storage    *utils.StorageClient
	userClient userPb.UserServiceClient
}

func NewBuildingUsecase(r repo.BuildingRepository, fr repo.FileRepository, s *utils.StorageClient, userClient userPb.UserServiceClient) BuildingUsecase {
	return BuildingUsecase{
		buildRepo:  r,
		fileRepo:   fr,
		storage:    s,
		userClient: userClient,
	}
}

func (u *BuildingUsecase) CreateBuilding(ctx context.Context, req *pb.CreateBuildingRequest) (*pb.BuildingResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	buildType := domain.BuildingTypeHouse
	if val, ok := helper.ProtoToDomainBuildingType[req.Type]; ok {
		buildType = val
	}

	var domainFiles []domain.FileAssociate
	if len(req.NewFiles) > 0 {
		for _, f := range req.NewFiles {
			domainFiles = append(domainFiles, domain.FileAssociate{
				Link:     f.Url,
				FileType: f.FileType,
				UserID:   userID,
			})
		}
	}

	var lands []domain.Land
	if len(req.LandIds) > 0 {
		for _, id := range req.LandIds {
			if uid, err := uuid.Parse(id); err == nil {
				lands = append(lands, domain.Land{ID: uid})
			}
		}
	}

	var ins []domain.Insurance
	if len(req.InsIds) > 0 {
		for _, id := range req.InsIds {
			if uid, err := uuid.Parse(id); err == nil {
				ins = append(ins, domain.Insurance{ID: uid})
			}
		}
	}

	loc := domain.Location{
		Address:     req.Location.Address,
		Subdistrict: req.Location.Subdistrict,
		District:    req.Location.District,
		Province:    req.Location.Province,
		PostalCode:  req.Location.PostalCode,
	}

	building := &domain.Building{
		UserID:      userID,
		Name:        req.Name,
		Type:        buildType,
		Area:        req.Area,
		Amount:      req.Amount,
		Description: req.Description,
		Location:    loc,
		Lands:       lands,
		Insurances:  ins,
		Files:       domainFiles,
	}

	if err := u.buildRepo.CreateBuilding(ctx, building); err != nil {
		return nil, err
	}

	return &pb.BuildingResponse{
		Building: mapper.ToBuildingProto(building),
	}, nil
}

func (u *BuildingUsecase) GetBuilding(ctx context.Context, req *pb.GetAssetRequest) (*pb.BuildingArrayResponse, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	buildings, err := u.buildRepo.GetBuilding(ctx, uid)
	if err != nil {
		return nil, err
	}

	var BuildList []*pb.Building
	for _, item := range buildings {
		BuildList = append(BuildList, mapper.ToBuildingProto(item))
	}

	return &pb.BuildingArrayResponse{
		Success:  true,
		Building: BuildList,
	}, nil
}

func (u *BuildingUsecase) GetBuildingByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.BuildingArrayResponse, error) {
	var ids []uuid.UUID
	for _, idStr := range req.Ids {
		if parsedID, err := uuid.Parse(idStr); err == nil {
			ids = append(ids, parsedID)
		}
	}

	bu, err := u.buildRepo.GetBuildingByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var pbBuildings []*pb.Building
	for _, a := range bu {
		pbBuildings = append(pbBuildings, mapper.ToBuildingProto(a))
	}

	return &pb.BuildingArrayResponse{
		Building: pbBuildings,
	}, nil
}

func (u *BuildingUsecase) GetBatchBuildingByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.BuildingArrayResponse, error) {
	var ids []uuid.UUID
	for _, idStr := range req.Ids {
		if parsedID, err := uuid.Parse(idStr); err == nil {
			ids = append(ids, parsedID)
		}
	}

	bu, err := u.buildRepo.GetBatchBuildingByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var pbBuildings []*pb.Building
	for _, a := range bu {
		pbBuildings = append(pbBuildings, mapper.ToBuildingProto(a))
	}

	return &pb.BuildingArrayResponse{
		Building: pbBuildings,
	}, nil
}

func (u *BuildingUsecase) GetBuildingByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.BuildingResponse, error) {
	id, err := uuid.Parse(req.Id)
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

	var addLandIDs, removeLandIDs, addInsIDs, removeInsIDs []uuid.UUID
	if len(req.LandIds) > 0 {
		for _, idStr := range req.LandIds {
			if parsedID, err := uuid.Parse(idStr); err == nil {
				addLandIDs = append(addLandIDs, parsedID)
			}
		}
	}

	if len(req.DeleteLandIds) > 0 {
		for _, idStr := range req.DeleteLandIds {
			if parsedID, err := uuid.Parse(idStr); err == nil {
				removeLandIDs = append(removeLandIDs, parsedID)
			}
		}
	}

	if len(req.InsIds) > 0 {
		for _, idStr := range req.InsIds {
			if parsedID, err := uuid.Parse(idStr); err == nil {
				addInsIDs = append(addInsIDs, parsedID)
			}
		}
	}

	if len(req.DeleteInsIds) > 0 {
		for _, idStr := range req.DeleteInsIds {
			if parsedID, err := uuid.Parse(idStr); err == nil {
				removeInsIDs = append(removeInsIDs, parsedID)
			}
		}
	}

	syncParams := domain.FileSyncParams{
		UserID:        uid,
		EntityID:      id,
		EntityType:    "building",
		NewFiles:      req.NewFiles,
		DeleteFileIDs: req.DeleteFileIds,
	}

	err = helper.SyncEntityFiles(ctx, u.fileRepo, u.storage, syncParams)
	if err != nil {
		return nil, err
	}

	updatedBuilding, err := u.buildRepo.UpdateBuilding(ctx, building, addLandIDs, removeLandIDs, addInsIDs, removeInsIDs)
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

	_, err = u.buildRepo.GetBuildingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := u.buildRepo.SoftDeleteBuilding(ctx, id, uid); err != nil {
		return nil, err
	}

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
		return err
	}

	for _, b := range expiredBuilding {
		helper.CleanupAssetResource(
			ctx,
			b.ID,
			b.Files,
			u.storage,
			u.userClient,
			func(id uuid.UUID) error {
				return u.buildRepo.HardDeleteBuilding(ctx, id)
			},
		)
	}

	return nil
}
