package usecase

import (
	"context"
	"errors"
	repo "wealth-vault/user-service/internal/repository/interface"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils/mapper"

	"github.com/google/uuid"
)

type MessageUsecase struct {
	msgRepo  repo.MsgRepository
	itemRepo repo.ShareItemRepository
}

func NewMessageUsecase(m repo.MsgRepository, i repo.ShareItemRepository) MessageUsecase {
	return MessageUsecase{
		msgRepo:  m,
		itemRepo: i,
	}
}

func (u *MessageUsecase) GetGroupMessages(ctx context.Context, req *pb.GetGroupMessagesRequest) (*pb.GetGroupMessagesResponse, error) {
	groupUUID, _ := uuid.Parse(req.GroupId)
	userUUID, _ := uuid.Parse(req.UserId)
	isMember, _ := u.itemRepo.IsGroupMember(ctx, groupUUID, userUUID)
	if !isMember {
		return nil, errors.New("unauthorized")
	}

	msgs, err := u.msgRepo.GetGroupMessages(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}

	var pbMsgs []*pb.MessageDetail
	for _, m := range msgs {
		pbMsgs = append(pbMsgs, mapper.ToGroupMessageProto(m, req.UserId))
	}

	return &pb.GetGroupMessagesResponse{
		Messages: pbMsgs,
	}, nil
}

func (u *MessageUsecase) GetPrivateMessages(ctx context.Context, req *pb.GetPrivateMessagesRequest) (*pb.GetPrivateMessagesResponse, error) {
	msgs, err := u.msgRepo.GetPrivateMessages(ctx, req.UserId, req.FriendId)
	if err != nil {
		return nil, err
	}

	var pbMsgs []*pb.MessageDetail
	for _, m := range msgs {
		pbMsgs = append(pbMsgs, mapper.ToPrivateMessageProto(m, req.UserId))
	}

	return &pb.GetPrivateMessagesResponse{
		Messages: pbMsgs,
	}, nil
}
