package usecase

import (
	"context"
	"errors"
	"wealth-vault/asset-service/internal/domain"
	repo "wealth-vault/asset-service/internal/repository/interface"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils"
	helper "wealth-vault/asset-service/pkg/utils/helper"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AssetUsecase struct {
	assetRepo repo.AssetRepository
	fileRepo  repo.FileRepository
	storage   *utils.StorageClient
}

func NewAssetUsecase(r repo.AssetRepository, fr repo.FileRepository, s *utils.StorageClient) AssetUsecase {
	return AssetUsecase{
		assetRepo: r,
		fileRepo:  fr,
		storage:   s,
	}
}

func (u *AssetUsecase) CreateAsset(ctx context.Context, req *pb.CreateAssetRequest) (*pb.CreateAssetResponse, error) {
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

	asset := &domain.Asset{
		UserID:               userID,
		Name:                 req.Name,
		Amount:               req.Amount,
		Type:                 domain.AssetType(req.Type.String()),
		IsIncludedInNetWorth: &req.IsIncludedInNetWorth,
		Description:          req.Description,
		Files:                domainFiles,
	}

	jsonBytes, err := helper.MapProtoDetailToJSON(req.Detail)
	if err != nil {
		return nil, errors.New("invalid detail data")
	}

	asset.Details = datatypes.JSON(jsonBytes)
	if err := u.assetRepo.CreateAsset(ctx, asset); err != nil {
		return nil, err
	}

	return &pb.CreateAssetResponse{
		Success: true,
		Id:      asset.ID.String(),
	}, nil
}

func (u *AssetUsecase) GetAsset(ctx context.Context, req *pb.GetAssetRequest) (*pb.AssetArrayResponse, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	assets, err := u.assetRepo.GetAsset(ctx, uid)
	if err != nil {
		return nil, err
	}

	var AssetList []*pb.Asset
	for _, item := range assets {
		AssetList = append(AssetList, utils.ToAssetProto(item))
	}

	return &pb.AssetArrayResponse{
		Success: true,
		Asset:   AssetList,
	}, nil
}

func (u *AssetUsecase) GetAssetByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.AssetResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	asset, err := u.assetRepo.GetAssetByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	return &pb.AssetResponse{
		Success: true,
		Asset:   utils.ToAssetProto(asset),
	}, nil
}

func (u *AssetUsecase) UpdateAsset(ctx context.Context, req *pb.UpdateAssetRequest) (*pb.AssetResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	asset, err := u.assetRepo.GetAssetByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	updateMask, err := helper.ApplyUpdateAssetFields(req, asset)
	if err != nil {
		return nil, err
	}

	syncParams := domain.FileSyncParams{
		UserID:        uid,
		EntityID:      id,
		EntityType:    "asset",
		NewFiles:      req.NewFiles,
		DeleteFileIDs: req.DeleteFileIds,
	}

	err = helper.SyncEntityFiles(ctx, u.fileRepo, u.storage, syncParams)
	if err != nil {
		return nil, err
	}

	updatedAsset, err := u.assetRepo.UpdateAsset(ctx, asset, updateMask)
	if err != nil {
		return nil, err
	}

	return &pb.AssetResponse{
		Success: true,
		Asset:   utils.ToAssetProto(updatedAsset),
	}, nil
}

func (u *AssetUsecase) DeleteAsset(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	existingCash, err := u.assetRepo.GetAssetByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}

	if err := u.assetRepo.DeleteAsset(ctx, id, uid); err != nil {
		return nil, err
	}

	if len(existingCash.Files) > 0 {
		fileURLs := make([]string, len(existingCash.Files))
		for i, f := range existingCash.Files {
			fileURLs[i] = f.Link
		}

		helper.DeleteFilesAsync(u.storage, fileURLs)
	}

	return &pb.DeleteAssetResponse{
		Success: true,
	}, nil
}
