package usecase

import (
	"context"
	"wealth-vault/user-service/internal/domain"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
)

type ShareItemUsecase interface {
	ShareItem(ctx context.Context, req *pb.ShareItemRequest) (*pb.ShareItemResponse, error)
	GetSharedIteminGroup(ctx context.Context, req *pb.GetGroupItemsRequest) (*pb.GetGroupItemsResponse, error)
	GetSharedIteminFriend(ctx context.Context, req *pb.GetFriendItemRequest) (*pb.GetFriendItemsResponse, error)
	UnsharedIteminGroup(ctx context.Context, req *pb.UnshareItemRequest) (*pb.ShareItemResponse, error)
	UnsharedIteminFriend(ctx context.Context, req *pb.UnshareItemRequest) (*pb.ShareItemResponse, error)
	AddMemberToGroup(ctx context.Context, req *pb.AddMemberRequest) (*pb.ActionResponse, error)
	GrantAccess(ctx context.Context, req *pb.GrantAccessRequest) (*pb.ActionResponse, error)
	ProcessScheduledEmails(ctx context.Context) error
	DeleteAllReferencesByEntityID(ctx context.Context, req *pb.DeleteByEntityRequest) (*pb.DeleteByEntityResponse, error)
	BatchShareAssets(ctx context.Context, req domain.BatchShareRequest) error
	GetItemSharedTargets(ctx context.Context, req *pb.GetItemSharedTargetsRequest) (*pb.GetItemSharedTargetsResponse, error)
	GetSharedItemIDs(ctx context.Context, req *pb.GetSharedItemIDsRequest) (*pb.GetSharedItemIDsResponse, error)
	GetItemsSharedByFriend(ctx context.Context, req *pb.GetItemsSharedByFriendRequest) (*pb.GetItemsSharedByFriendResponse, error)
	GetAllSharedItemIDsByUser(ctx context.Context, req *pb.GetAllSharedItemIDsByUserRequest) (*pb.GetAllSharedItemIDsByUserResponse, error)
}
