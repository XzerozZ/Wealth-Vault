package usecase

import (
	"context"
	"errors"
	"wealth-vault/asset-service/internal/domain"
	repo "wealth-vault/asset-service/internal/repository/interface"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils"
	helper "wealth-vault/asset-service/pkg/utils/helper"
	"wealth-vault/asset-service/pkg/utils/mapper"

	"github.com/google/uuid"
)

type LandUsecase struct {
	landRepo repo.LandRepository
	fileRepo repo.FileRepository
	storage  *utils.StorageClient
}

func NewLandUsecase(r repo.LandRepository, fr repo.FileRepository, s *utils.StorageClient) LandUsecase {
	return LandUsecase{
		landRepo: r,
		fileRepo: fr,
		storage:  s,
	}
}

func (u *LandUsecase) CreateLand(ctx context.Context, req *pb.CreateLandRequest) (*pb.LandResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
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

	var builds []domain.Building
	if len(req.BuildingIds) > 0 {
		for _, id := range req.BuildingIds {
			if uid, err := uuid.Parse(id); err == nil {
				builds = append(builds, domain.Building{ID: uid})
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

	land := &domain.Land{
		UserID:      userID,
		Name:        req.Name,
		DeedNum:     req.DeedNum,
		Area:        req.Area,
		Amount:      req.Amount,
		Description: req.Description,
		Location:    loc,
		Buildings:   builds,
		Files:       domainFiles,
	}

	if err := u.landRepo.CreateLand(ctx, land); err != nil {
		return nil, err
	}

	return &pb.LandResponse{
		Land: mapper.ToLandProto(land),
	}, nil
}

func (u *LandUsecase) GetLand(ctx context.Context, req *pb.GetAssetRequest) (*pb.LandArrayResponse, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	lands, err := u.landRepo.GetLand(ctx, uid)
	if err != nil {
		return nil, err
	}

	var LandList []*pb.Land
	for _, item := range lands {
		LandList = append(LandList, mapper.ToLandProto(item))
	}

	return &pb.LandArrayResponse{
		Success: true,
		Land:    LandList,
	}, nil
}

func (u *LandUsecase) GetLandByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.LandArrayResponse, error) {
	var ids []uuid.UUID
	for _, idStr := range req.Ids {
		if parsedID, err := uuid.Parse(idStr); err == nil {
			ids = append(ids, parsedID)
		}
	}

	bu, err := u.landRepo.GetLandByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var pbLand []*pb.Land
	for _, a := range bu {
		pbLand = append(pbLand, mapper.ToLandProto(a))
	}

	return &pb.LandArrayResponse{
		Land: pbLand,
	}, nil
}

func (u *LandUsecase) GetLandByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.LandResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	land, err := u.landRepo.GetLandByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	return &pb.LandResponse{
		Success: true,
		Land:    mapper.ToLandProto(land),
	}, nil
}

func (u *LandUsecase) UpdateLand(ctx context.Context, req *pb.UpdateLandRequest) (*pb.LandResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.Land.UserId)
	if err != nil {
		return nil, err
	}

	land, err := u.landRepo.GetLandByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	if err := helper.ApplyUpdateLandFields(req, land); err != nil {
		return nil, err
	}

	var addBuildIDs, removeBuildIDs []uuid.UUID

	if len(req.BuildingIds) > 0 {
		for _, idStr := range req.BuildingIds {
			if parsedID, err := uuid.Parse(idStr); err == nil {
				addBuildIDs = append(addBuildIDs, parsedID)
			}
		}
	}

	if len(req.DeleteBuildingIds) > 0 {
		for _, idStr := range req.DeleteBuildingIds {
			if parsedID, err := uuid.Parse(idStr); err == nil {
				removeBuildIDs = append(removeBuildIDs, parsedID)
			}
		}
	}

	syncParams := domain.FileSyncParams{
		UserID:        uid,
		EntityID:      id,
		EntityType:    "land",
		NewFiles:      req.NewFiles,
		DeleteFileIDs: req.DeleteFileIds,
	}

	err = helper.SyncEntityFiles(ctx, u.fileRepo, u.storage, syncParams)
	if err != nil {
		return nil, err
	}

	updatedLand, err := u.landRepo.UpdateLand(ctx, land, addBuildIDs, removeBuildIDs)
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

	existingLand, err := u.landRepo.GetLandByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	if err := u.landRepo.DeleteLand(ctx, id, uid); err != nil {
		return nil, err
	}

	if len(existingLand.Files) > 0 {
		fileURLs := make([]string, len(existingLand.Files))
		for i, f := range existingLand.Files {
			fileURLs[i] = f.Link
		}

		helper.DeleteFilesAsync(u.storage, fileURLs)
	}

	return &pb.DeleteAssetResponse{
		Success: true,
	}, nil
}
