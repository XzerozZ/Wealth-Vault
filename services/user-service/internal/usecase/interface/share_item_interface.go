package usecase

import (
	"context"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
)

type GroupItemUsecase interface {
	ShareItemtoGroup(ctx context.Context, req *pb.ShareItemRequest) (*pb.ShareItemResponse, error)
	GetSharedIteminGroup(ctx context.Context, req *pb.GetGroupItemsRequest) (*pb.GetGroupItemsResponse, error)
	GetSharedIteminFriend(ctx context.Context, req *pb.GetFriendItemRequest) (*pb.GetFriendItemsResponse, error)
	UnsharedIteminGroup(ctx context.Context, req *pb.UnshareItemRequest) (*pb.ShareItemResponse, error)
	UnsharedIteminFriend(ctx context.Context, req *pb.UnshareItemRequest) (*pb.ShareItemResponse, error)
	AddMemberToGroup(ctx context.Context, req *pb.AddMemberRequest) (*pb.ActionResponse, error)
	GrantAccess(ctx context.Context, req *pb.GrantAccessRequest) (*pb.ActionResponse, error)
	ProcessScheduledEmails(ctx context.Context) error
}
