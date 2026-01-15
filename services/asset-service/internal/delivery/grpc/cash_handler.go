package grpc

import (
	"context"
	"wealth-vault/asset-service/internal/domain"
	"wealth-vault/asset-service/internal/usecase"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AssetGRPCHandler struct {
	pb.UnimplementedAssetServiceServer
	usecase usecase.AssetUsecase
}

func NewAssetGRPCHandler(u usecase.AssetUsecase) *AssetGRPCHandler {
	return &AssetGRPCHandler{usecase: u}
}

func (h *AssetGRPCHandler) CreateCash(ctx context.Context, req *pb.CreateCashRequest) (*pb.CreateCashResponse, error) {
	id, err := h.usecase.CreateCash(ctx, &domain.Cash{
		Name:        req.Name,
		Value:       req.Amount,
		Description: req.Description,
		UserID:      req.CreatedBy,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateCashResponse{
		Success: true,
		CashId:  id,
	}, nil
}

func (h *AssetGRPCHandler) GetCash(ctx context.Context, req *pb.GetCashRequest) (*pb.CashArrayResponse, error) {
	cash, err := h.usecase.GetCash(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	var CashList []*pb.Cash
	for _, item := range cash {
		CashList = append(CashList, &pb.Cash{
			Id:          item.ID,
			Name:        item.Name,
			Amount:      item.Value,
			Description: item.Description,
			CreatedBy:   item.UserID,
			CreatedAt:   timestamppb.New(item.CreatedAt),
			UpdatedAt:   timestamppb.New(item.UpdatedAt),
		})
	}

	return &pb.CashArrayResponse{
		Success: true,
		Cash:    CashList,
	}, nil
}

func (h *AssetGRPCHandler) GetCashByID(ctx context.Context, req *pb.GetCashByIDRequest) (*pb.CashResponse, error) {
	cash, err := h.usecase.GetCashByID(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.CashResponse{
		Success: true,
		Cash: &pb.Cash{
			Id:          cash.ID,
			Name:        cash.Name,
			Amount:      cash.Value,
			Description: cash.Description,
			CreatedBy:   cash.UserID,
			CreatedAt:   timestamppb.New(cash.CreatedAt),
			UpdatedAt:   timestamppb.New(cash.UpdatedAt),
		},
	}, nil
}

func (h *AssetGRPCHandler) UpdateCash(ctx context.Context, req *pb.UpdateCashRequest) (*pb.CashResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	reqCash := req.GetCash()
	if reqCash == nil {
		return nil, status.Error(codes.InvalidArgument, "user data is required")
	}

	var mask []string
	if req.GetUpdateMask() != nil {
		mask = req.GetUpdateMask().GetPaths()
	}

	input := &domain.UpdateCashInput{
		ID:          req.Id,
		Name:        reqCash.GetName(),
		Value:       reqCash.GetAmount(),
		Description: reqCash.GetDescription(),
		UserID:      reqCash.GetCreatedBy(),
		UpdateMask:  mask,
	}

	cash, err := h.usecase.UpdateCash(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CashResponse{
		Success: true,
		Cash: &pb.Cash{
			Id:          cash.ID,
			Name:        cash.Name,
			Amount:      cash.Value,
			Description: cash.Description,
			CreatedBy:   cash.UserID,
			CreatedAt:   timestamppb.New(cash.CreatedAt),
			UpdatedAt:   timestamppb.New(cash.UpdatedAt),
		},
	}, nil
}
