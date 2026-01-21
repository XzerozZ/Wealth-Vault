package usecase

import (
	"context"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	GetUser(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error)
}
