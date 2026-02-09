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

type LandUsecase struct {
	landRepo   repo.LandRepository
	fileRepo   repo.FileRepository
	storage    *utils.StorageClient
	userClient userPb.UserServiceClient
}

func NewLandUsecase(r repo.LandRepository, fr repo.FileRepository, s *utils.StorageClient, userClient userPb.UserServiceClient) LandUsecase {
	return LandUsecase{
		landRepo:   r,
		fileRepo:   fr,
		storage:    s,
		userClient: userClient,
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

func (u *LandUsecase) GetBatchLandByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.LandArrayResponse, error) {
	var ids []uuid.UUID
	for _, idStr := range req.Ids {
		if parsedID, err := uuid.Parse(idStr); err == nil {
			ids = append(ids, parsedID)
		}
	}

	bu, err := u.landRepo.GetBatchLandByIDs(ctx, ids)
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
	id, err := uuid.Parse(req.Id)
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

	_, err = u.landRepo.GetLandByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := u.landRepo.SoftDeleteLand(ctx, id, uid); err != nil {
		return nil, err
	}

	return &pb.DeleteAssetResponse{
		Success: true,
	}, nil
}

func (u *LandUsecase) CleanupExpiredLand(ctx context.Context) error {
	cutoffTime := time.Now().AddDate(0, 0, -7)
	GetExpiredLand, err := u.landRepo.GetExpiredLand(ctx, cutoffTime)
	if err != nil {
		return err
	}

	if len(GetExpiredLand) == 0 {
		return err
	}

	for _, l := range GetExpiredLand {
		helper.CleanupAssetResource(
			ctx,
			l.ID,
			l.Files,
			u.storage,
			u.userClient,
			func(id uuid.UUID) error {
				return u.landRepo.HardDeleteLand(ctx, id)
			},
		)
	}

	return nil
}
