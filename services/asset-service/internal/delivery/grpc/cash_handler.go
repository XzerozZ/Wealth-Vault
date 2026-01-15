package grpc

import (
	"context"
	"wealth-vault/asset-service/internal/domain"
	"wealth-vault/asset-service/internal/usecase"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CashGRPCHandler struct {
	pb.UnimplementedAssetServiceServer
	usecase usecase.CashUsecase
}

func NewCashGRPCHandler(u usecase.CashUsecase) *CashGRPCHandler {
	return &CashGRPCHandler{usecase: u}
}

func (h *CashGRPCHandler) CreateCash(ctx context.Context, req *pb.CreateCashRequest) (*pb.CreateCashResponse, error) {
	var domainFiles []domain.FileAssociate
	if len(req.Files) > 0 {
		for _, f := range req.Files {
			domainFiles = append(domainFiles, domain.FileAssociate{
				Link:     f.Url,
				FileType: f.FileType,
				UserID:   req.CreatedBy,
			})
		}
	}

	id, err := h.usecase.CreateCash(ctx, &domain.Cash{
		Name:        req.Name,
		Value:       req.Amount,
		Description: req.Description,
		UserID:      req.CreatedBy,
		Files:       domainFiles,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateCashResponse{
		Success: true,
		CashId:  id,
	}, nil
}

func (h *CashGRPCHandler) GetCash(ctx context.Context, req *pb.GetCashRequest) (*pb.CashArrayResponse, error) {
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

func (h *CashGRPCHandler) GetCashByID(ctx context.Context, req *pb.CashByIDRequest) (*pb.CashResponse, error) {
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
			Files:       utils.ToPbFiles(cash.Files),
		},
	}, nil
}

func (h *CashGRPCHandler) UpdateCash(ctx context.Context, req *pb.UpdateCashRequest) (*pb.CashResponse, error) {
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

	var newFiles []domain.FileAssociate
	if len(req.NewFiles) > 0 {
		for _, f := range req.NewFiles {
			newFiles = append(newFiles, domain.FileAssociate{
				Link:     f.Url,
				FileType: f.FileType,
				UserID:   reqCash.GetCreatedBy(),
			})
		}
	}

	input := &domain.UpdateCashInput{
		ID:            req.Id,
		Name:          reqCash.GetName(),
		Value:         reqCash.GetAmount(),
		Description:   reqCash.GetDescription(),
		UserID:        reqCash.GetCreatedBy(),
		UpdateMask:    mask,
		NewFiles:      newFiles,
		DeleteFileIDs: req.DeleteFileIds,
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
			Files:       utils.ToPbFiles(cash.Files),
		},
	}, nil
}

func (h *CashGRPCHandler) DeleteCash(ctx context.Context, req *pb.CashByIDRequest) (*pb.CashResponse, error) {
	err := h.usecase.DeleteCash(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.CashResponse{
		Success: true,
	}, nil
}
