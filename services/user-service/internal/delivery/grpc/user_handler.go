package grpc

import (
	"context"
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/usecase"
	userpb "wealth-vault/user-service/pkg/pb/proto/user"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		Id: id.String(),
	}, nil
}

func (h *UserGRPCHandler) GetUser(ctx context.Context, req *userpb.GetUserByIDRequest) (*userpb.UserResponse, error) {
	user, err := h.usecase.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{
		Success: true,
		User: &userpb.User{
			Id:          user.ID.String(),
			Email:       user.Email,
			Firstname:   user.Firstname,
			Lastname:    user.Lastname,
			Username:    user.Username,
			Profile:     user.Profile,
			Phonenumber: user.Phonenumber,
			Birthday:    user.Birthday.Format("2006-01-02"),
			CreatedAt:   timestamppb.New(user.CreatedAt),
			UpdatedAt:   timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

func (h *UserGRPCHandler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
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

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid format")
	}

	input := &domain.UpdateUserInput{
		ID:          id,
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

	return &userpb.UserResponse{
		Success: true,
		User: &userpb.User{
			Id:          user.ID.String(),
			Email:       user.Email,
			Firstname:   user.Firstname,
			Lastname:    user.Lastname,
			Username:    user.Username,
			Profile:     user.Profile,
			Phonenumber: user.Phonenumber,
			Birthday:    user.Birthday.Format("2006-01-02"),
			CreatedAt:   timestamppb.New(user.CreatedAt),
			UpdatedAt:   timestamppb.New(user.UpdatedAt),
		},
	}, nil
}
