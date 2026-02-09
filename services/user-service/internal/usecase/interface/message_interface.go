package usecase

import (
	"context"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
)

type MessageUsecase interface {
	GetGroupMessages(ctx context.Context, req *pb.GetGroupMessagesRequest) (*pb.GetGroupMessagesResponse, error)
	GetPrivateMessages(ctx context.Context, req *pb.GetPrivateMessagesRequest) (*pb.GetPrivateMessagesResponse, error)
}
