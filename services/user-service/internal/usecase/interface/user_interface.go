package usecase

import (
	"context"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error)
	GetUser(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error)
	GetUsersByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserInfoResponse, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error)
	GetFriendList(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.FriendListResponse, error)
	GetPendingRequests(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.FriendListResponse, error)
	AddFriend(ctx context.Context, req *pb.FriendRequest) (*pb.FriendResponse, error)
	VerifyFriendship(ctx context.Context, req *pb.CheckFriendshipRequest) (*pb.CheckFriendshipResponse, error)
	AcceptFriend(ctx context.Context, req *pb.AcceptFriendRequest) (*pb.FriendResponse, error)
	SetCloseFriend(ctx context.Context, req *pb.SetCloseFriendRequest) (*pb.SetCloseFriendResponse, error)
	GetCloseFriends(ctx context.Context, req *pb.GetCloseFriendsRequest) (*pb.GetCloseFriendsResponse, error)
	DeleteFriend(ctx context.Context, req *pb.FriendRequest) (*pb.FriendResponse, error)
	ProcessLegacyAutoShare(ctx context.Context) error
}
