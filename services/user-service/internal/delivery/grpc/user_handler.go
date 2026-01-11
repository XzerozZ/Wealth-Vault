package grpc

import (
	"context"
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/usecase"
	userpb "wealth-vault/user-service/pkg/pb/proto/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (h *UserGRPCHandler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	reqUser := req.GetUser()
	if reqUser == nil {
		return nil, status.Error(codes.InvalidArgument, "user data is required")
	}

	var mask []string
	if req.GetUpdateMask() != nil {
		mask = req.GetUpdateMask().GetPaths()
	}

	input := &domain.UpdateUserInput{
		ID:          req.Id,
		Firstname:   reqUser.GetFirstname(),
		Lastname:    reqUser.GetLastname(),
		Username:    reqUser.GetUsername(),
		Profile:     reqUser.GetProfile(),
		Phonenumber: reqUser.GetPhonenumber(),
		BirthdayStr: reqUser.GetBirthday(),
		UpdateMask:  mask,
	}

	user, err := h.usecase.UpdateUser(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.UpdateUserResponse{
		Success: true,
		User: &userpb.User{
			Id:          user.ID,
			Email:       user.Email,
			Firstname:   user.Firstname,
			Lastname:    user.Lastname,
			Username:    user.Username,
			Profile:     user.Profile,
			Phonenumber: user.Phonenumber,
			Birthday:    user.Birthday.Format("2006-01-02"),
		},
	}, nil
}
