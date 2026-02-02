package usecase

import (
	"context"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

type InvestmentUsecase interface {
	CreateInvestment(ctx context.Context, req *pb.CreateInvestmentRequest) (*pb.InvestmentResponse, error)
	GetInvestment(ctx context.Context, req *pb.GetAssetRequest) (*pb.InvestmentArrayResponse, error)
	GetInvestmentByIDs(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.InvestmentArrayResponse, error)
	GetInvestmentByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.InvestmentResponse, error)
	UpdateInvestment(ctx context.Context, req *pb.UpdateInvestmentRequest) (*pb.InvestmentResponse, error)
	DeleteInvestment(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error)
}
