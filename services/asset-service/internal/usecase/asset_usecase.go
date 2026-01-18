package usecase

import (
	"context"
	"errors"
	"fmt"
	"wealth-vault/asset-service/internal/domain"
	repo "wealth-vault/asset-service/internal/repository/interface"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils"

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
		IsIncludedInNetWorth: req.IsIncludedInNetWorth,
		Description:          req.Description,
		Files:                domainFiles,
	}

	jsonBytes, err := utils.MapProtoDetailToJSON(req.Detail)
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
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, errors.New("invalid asset id")
	}

	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
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
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, errors.New("invalid asset id")
	}

	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	asset, err := u.assetRepo.GetAssetByID(ctx, id, uid)
	if err != nil {
		return nil, err
	}
	fmt.Println(req)
	var updateMask []string
	has := func(target string) bool {
		if req.UpdateMask == nil || len(req.UpdateMask.Paths) == 0 {
			return true
		}

		for _, p := range req.UpdateMask.Paths {
			if p == target {
				return true
			}
		}

		return false
	}

	if has("name") {
		asset.Name = req.Name
		updateMask = append(updateMask, "Name")
	}
	if has("amount") {
		asset.Amount = req.Amount
		updateMask = append(updateMask, "Amount")
	}
	if has("description") {
		asset.Description = req.Description
		updateMask = append(updateMask, "Description")
	}
	if has("detail") && req.Detail != nil {
		jsonBytes, err := utils.MapProtoDetailToJSON(req.Detail)
		if err != nil {
			return nil, err
		}

		asset.Details = datatypes.JSON(jsonBytes)
		updateMask = append(updateMask, "Details")
	}

	if len(req.DeleteFileIds) > 0 {
		var fileUUIDs []uuid.UUID
		for _, idStr := range req.DeleteFileIds {
			parsedID, err := uuid.Parse(idStr)
			if err == nil {
				fileUUIDs = append(fileUUIDs, parsedID)
			}
		}

		filesToDelete, err := u.fileRepo.GetFilesByIDs(ctx, fileUUIDs)
		if err == nil {
			var validURLs []string
			for _, f := range filesToDelete {
				if f.UserID == uid {
					validURLs = append(validURLs, f.Link)
				}
			}

			if len(validURLs) > 0 {
				utils.DeleteFilesAsync(u.storage, validURLs)
			}
		}

		err = u.fileRepo.DeleteFiles(ctx, req.DeleteFileIds, id, uid)
		if err != nil {
			return nil, err
		}
	}

	if len(req.NewFiles) > 0 {
		var filesToCreate []domain.FileAssociate

		for _, f := range req.NewFiles {
			filesToCreate = append(filesToCreate, domain.FileAssociate{
				ID:         uuid.New(),
				EntityID:   id,
				EntityType: "asset",
				UserID:     asset.UserID,
				Link:       f.Url,
				FileType:   f.FileType,
			})
		}

		if err := u.fileRepo.CreateFiles(ctx, filesToCreate); err != nil {
			return nil, err
		}
	}

	updatedAsset, err := u.assetRepo.UpdateAsset(ctx, asset, updateMask)
	if err != nil {
		return nil, err
	}

	fmt.Println(updatedAsset)
	return &pb.AssetResponse{
		Success: true,
		Asset:   utils.ToAssetProto(updatedAsset),
	}, nil
}

func (u *AssetUsecase) DeleteAsset(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, errors.New("invalid asset id")
	}

	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
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

		utils.DeleteFilesAsync(u.storage, fileURLs)
	}

	return &pb.DeleteAssetResponse{
		Success: true,
	}, nil
}
