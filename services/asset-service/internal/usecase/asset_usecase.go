package usecase

import (
	"context"
	repo "wealth-vault/asset-service/internal/repository/interface"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils"
)

type AssetUsecase struct {
	assetRepo repo.AssetRepository
}

func NewAssetUsecase(r repo.AssetRepository) AssetUsecase {
	return AssetUsecase{assetRepo: r}
}

func (u *AssetUsecase) CheckExists(ctx context.Context, req *pb.CheckAssetRequest) (*pb.CheckAssetResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	res, err := u.assetRepo.CheckExists(ctx, req.Type, id, uid)
	if err != nil {
		return nil, err
	}

	return &pb.CheckAssetResponse{
		Exists: res,
	}, nil
}
