package usecase

import (
	"context"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

type LiabilityUsecase interface {
	CreateLiability(ctx context.Context, req *pb.CreateLiabilityRequest) (*pb.CreateLiabilityResponse, error)
	GetLiability(ctx context.Context, req *pb.GetLiabilityRequest) (*pb.LiabilityArrayResponse, error)
	GetLiabilityByID(ctx context.Context, req *pb.GetLiabilityByIDRequest) (*pb.LiabilityResponse, error)
	UpdateLiability(ctx context.Context, req *pb.UpdateLiabilityRequest) (*pb.LiabilityResponse, error)
	DeleteLiability(ctx context.Context, req *pb.DeleteLiabilityRequest) (*pb.DeleteLiabilityResponse, error)
}
