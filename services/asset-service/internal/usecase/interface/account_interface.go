package usecase

import (
	"context"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

type AccountUsecase interface {
	CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.AccountResponse, error)
	GetAccount(ctx context.Context, req *pb.GetAssetRequest) (*pb.AccountArrayResponse, error)
	GetAccountByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.AccountResponse, error)
	UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.AccountResponse, error)
	DeleteAccount(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error)
}
