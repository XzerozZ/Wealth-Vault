package usecase

import (
	"context"
	"errors"
	"wealth-vault/user-service/internal/domain"
	repo "wealth-vault/user-service/internal/repository/interface"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils"
	helper "wealth-vault/user-service/pkg/utils/helper"

	"github.com/google/uuid"
)

type UserUsecase struct {
	userRepo repo.UserRepository
	storage  *utils.StorageClient
}

func NewUserUsecase(r repo.UserRepository, s *utils.StorageClient) UserUsecase {
	return UserUsecase{
		userRepo: r,
		storage:  s,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &domain.User{
		Email:    req.Email,
		Username: req.Username,
	}

	if err := u.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		Id: user.ID.String(),
	}, nil
}

func (u *UserUsecase) GetUser(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := u.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Success: true,
		User:    utils.ToUserProto(user),
	}, nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := u.userRepo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	updateMask, err := helper.ApplyUpdateUserFields(req, u.storage, user)
	if err != nil {
		return nil, err
	}

	updatedUser, err := u.userRepo.UpdateUser(ctx, user, updateMask)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Success: true,
		User:    utils.ToUserProto(updatedUser),
	}, nil
}
