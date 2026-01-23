package usecase

import (
	"context"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error)
	GetUser(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error)
	GetFriendList(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.FriendListResponse, error)
	AddFriend(ctx context.Context, req *pb.FriendRequest) (*pb.FriendResponse, error)
}
