package grpc

import (
	"context"
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/usecase"
	userpb "wealth-vault/user-service/proto/userpb"
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
		Firstname:   req.Firstname,
		Lastname:    req.Lastname,
		Username:    req.Username,
		Phonenumber: req.Phonenumber,
	})
	if err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponse{
		Id: id,
	}, nil
}

func (h *UserGRPCHandler) GetUserByID(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, err := h.usecase.GetByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserResponse{
		Id:          user.ID,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Username:    user.Username,
		Phonenumber: user.Phonenumber,
	}, nil
}
