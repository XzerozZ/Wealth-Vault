package grpc

import (
	"context"
	"wealth-vault/user-service/internal/usecase"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
)

type UserGRPCHandler struct {
	pb.UnimplementedUserServiceServer
	usecase usecase.UserUsecase
}

func NewUserGRPCHandler(u usecase.UserUsecase) *UserGRPCHandler {
	return &UserGRPCHandler{usecase: u}
}

func (h *UserGRPCHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	res, err := h.usecase.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) GetUser(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error) {
	res, err := h.usecase.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *UserGRPCHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	res, err := h.usecase.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
