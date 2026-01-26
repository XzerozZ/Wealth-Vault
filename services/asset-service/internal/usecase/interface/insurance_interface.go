package usecase

import (
	"context"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

type InsuranceUsecase interface {
	CreateInsurance(ctx context.Context, req *pb.CreateInsuranceRequest) (*pb.InsuranceResponse, error)
	GetInsurance(ctx context.Context, req *pb.GetAssetRequest) (*pb.InsuranceArrayResponse, error)
	GetInsuranceByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.InsuranceResponse, error)
	UpdateInsurance(ctx context.Context, req *pb.UpdateInsuranceRequest) (*pb.InsuranceResponse, error)
	DeleteInsurance(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error)
}
