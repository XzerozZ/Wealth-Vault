package usecase

import (
	"context"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

type LandUsecase interface {
	CreateLand(ctx context.Context, req *pb.CreateLandRequest) (*pb.LandResponse, error)
	GetLand(ctx context.Context, req *pb.GetAssetRequest) (*pb.LandArrayResponse, error)
	GetLandByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.LandArrayResponse, error)
	GetBatchLandByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.LandArrayResponse, error)
	GetLandByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.LandResponse, error)
	UpdateLand(ctx context.Context, req *pb.UpdateLandRequest) (*pb.LandResponse, error)
	DeleteLand(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error)
	CleanupExpiredLand(ctx context.Context) error
}
