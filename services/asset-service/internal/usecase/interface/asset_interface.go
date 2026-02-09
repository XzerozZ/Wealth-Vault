package usecase

import (
	"context"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

type AssetUsecase interface {
	CheckExists(ctx context.Context, req *pb.CheckAssetRequest) (*pb.CheckAssetResponse, error)
	GetAllAssetIDs(ctx context.Context, req *pb.GetMyAssetsRequest) (*pb.GetMyAssetsResponse, error)
}
