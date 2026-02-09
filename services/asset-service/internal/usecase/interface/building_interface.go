package usecase

import (
	"context"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

type BuildingUsecase interface {
	CreateBuilding(ctx context.Context, req *pb.CreateBuildingRequest) (*pb.BuildingResponse, error)
	GetBuilding(ctx context.Context, req *pb.GetAssetRequest) (*pb.BuildingArrayResponse, error)
	GetBuildingByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.BuildingArrayResponse, error)
	GetBatchBuildingByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.BuildingArrayResponse, error)
	GetBuildingByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.BuildingResponse, error)
	UpdateBuilding(ctx context.Context, req *pb.UpdateBuildingRequest) (*pb.BuildingResponse, error)
	DeleteBuilding(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error)
	CleanupExpiredBuildings(ctx context.Context) error
}
