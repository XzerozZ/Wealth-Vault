package grpc

import (
	"context"
	"wealth-vault/asset-service/internal/usecase"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

type AssetGRPCHandler struct {
	pb.UnimplementedAssetServiceServer
	usecase usecase.AssetUsecase
}

func NewAssetGRPCHandler(u usecase.AssetUsecase) *AssetGRPCHandler {
	return &AssetGRPCHandler{usecase: u}
}

func (h *AssetGRPCHandler) CreateAsset(ctx context.Context, req *pb.CreateAssetRequest) (*pb.CreateAssetResponse, error) {
	res, err := h.usecase.CreateAsset(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetAsset(ctx context.Context, req *pb.GetAssetRequest) (*pb.AssetArrayResponse, error) {
	assets, err := h.usecase.GetAsset(ctx, req)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (h *AssetGRPCHandler) GetAssetByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.AssetResponse, error) {
	asset, err := h.usecase.GetAssetByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (h *AssetGRPCHandler) UpdateAsset(ctx context.Context, req *pb.UpdateAssetRequest) (*pb.AssetResponse, error) {
	asset, err := h.usecase.UpdateAsset(ctx, req)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (h *AssetGRPCHandler) DeleteAsset(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	asset, err := h.usecase.DeleteAsset(ctx, req)
	if err != nil {
		return nil, err
	}

	return asset, nil
}
