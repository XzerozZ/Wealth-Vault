package grpc

import (
	"context"
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/usecase"
	userpb "wealth-vault/user-service/pkg/pb/proto/user"
)

type UserGRPCHandler struct {
	userpb.UnimplementedUserServiceServer
	usecase usecase.UserUsecase
}

func NewUserGRPCHandler(u usecase.UserUsecase) *UserGRPCHandler {
	return &UserGRPCHandler{usecase: u}
}

func (h *UserGRPCHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	id, err := h.usecase.CreateUser(ctx, &domain.User{
		Email:    req.Email,
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponse{
		Id: id,
	}, nil
}
