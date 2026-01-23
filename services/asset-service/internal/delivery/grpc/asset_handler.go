package grpc

import (
	"context"
	"wealth-vault/asset-service/internal/usecase"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

type AssetGRPCHandler struct {
	pb.UnimplementedAssetServiceServer
	assetusecase usecase.AssetUsecase
	liausecase   usecase.LiabilityUsecase
}

func NewAssetGRPCHandler(au usecase.AssetUsecase, lu usecase.LiabilityUsecase) *AssetGRPCHandler {
	return &AssetGRPCHandler{
		assetusecase: au,
		liausecase:   lu,
	}
}

func (h *AssetGRPCHandler) CreateAsset(ctx context.Context, req *pb.CreateAssetRequest) (*pb.AssetResponse, error) {
	res, err := h.assetusecase.CreateAsset(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) CreateLiability(ctx context.Context, req *pb.CreateLiabilityRequest) (*pb.LiabilityResponse, error) {
	res, err := h.liausecase.CreateLiability(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetAsset(ctx context.Context, req *pb.GetAssetRequest) (*pb.AssetArrayResponse, error) {
	assets, err := h.assetusecase.GetAsset(ctx, req)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (h *AssetGRPCHandler) GetAssetByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.AssetResponse, error) {
	asset, err := h.assetusecase.GetAssetByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (h *AssetGRPCHandler) GetLiability(ctx context.Context, req *pb.GetLiabilityRequest) (*pb.LiabilityArrayResponse, error) {
	lias, err := h.liausecase.GetLiability(ctx, req)
	if err != nil {
		return nil, err
	}

	return lias, nil
}

func (h *AssetGRPCHandler) GetLiabilityByID(ctx context.Context, req *pb.GetLiabilityByIDRequest) (*pb.LiabilityResponse, error) {
	lia, err := h.liausecase.GetLiabilityByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return lia, nil
}

func (h *AssetGRPCHandler) UpdateAsset(ctx context.Context, req *pb.UpdateAssetRequest) (*pb.AssetResponse, error) {
	asset, err := h.assetusecase.UpdateAsset(ctx, req)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (h *AssetGRPCHandler) UpdateLiability(ctx context.Context, req *pb.UpdateLiabilityRequest) (*pb.LiabilityResponse, error) {
	lia, err := h.liausecase.UpdateLiability(ctx, req)
	if err != nil {
		return nil, err
	}

	return lia, nil
}

func (h *AssetGRPCHandler) DeleteAsset(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	asset, err := h.assetusecase.DeleteAsset(ctx, req)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (h *AssetGRPCHandler) DeleteLiability(ctx context.Context, req *pb.DeleteLiabilityRequest) (*pb.DeleteLiabilityResponse, error) {
	lia, err := h.liausecase.DeleteLiability(ctx, req)
	if err != nil {
		return nil, err
	}

	return lia, nil
}
