package usecase

import (
	"context"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

type CashUsecase interface {
	CreateCash(ctx context.Context, req *pb.CreateCashRequest) (*pb.CashResponse, error)
	GetCash(ctx context.Context, req *pb.GetAssetRequest) (*pb.CashArrayResponse, error)
	GetCashByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.CashArrayResponse, error)
	GetBatchCashByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.CashArrayResponse, error)
	GetCashByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.CashResponse, error)
	UpdateCash(ctx context.Context, req *pb.UpdateCashRequest) (*pb.CashResponse, error)
	DeleteCash(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error)
	CleanupExpiredCashes(ctx context.Context) error
}
