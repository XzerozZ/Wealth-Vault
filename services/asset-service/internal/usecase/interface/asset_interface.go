package usecase

import (
	"context"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

type AssetUsecase interface {
	CreateAsset(ctx context.Context, req *pb.CreateAssetRequest) (*pb.AssetResponse, error)
	GetAsset(ctx context.Context, req *pb.GetAssetRequest) (*pb.AssetArrayResponse, error)
	GetAssetByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.AssetResponse, error)
	UpdateAsset(ctx context.Context, req *pb.UpdateAssetRequest) (*pb.AssetResponse, error)
	DeleteAsset(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error)
}
