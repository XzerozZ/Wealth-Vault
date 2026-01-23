package usecase

import (
	"context"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
)

type GroupUsecase interface {
	CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.GroupResponse, error)
	GetMember(ctx context.Context, req *pb.GetGroupMembersRequest) (*pb.GetGroupMembersResponse, error)
	GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.GroupResponse, error)
	UpdateGroup(ctx context.Context, req *pb.UpdateGroupRequest) (*pb.GroupResponse, error)
}
