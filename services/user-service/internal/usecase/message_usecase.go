package usecase

import (
	"context"
	"errors"
	repo "wealth-vault/user-service/internal/repository/interface"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils"
	"wealth-vault/user-service/pkg/utils/mapper"
)

type MessageUsecase struct {
	msgRepo  repo.MsgRepository
	itemRepo repo.ShareItemRepository
}

func NewMessageUsecase(m repo.MsgRepository, i repo.ShareItemRepository) *MessageUsecase {
	return &MessageUsecase{
		msgRepo:  m,
		itemRepo: i,
	}
}

func (u *MessageUsecase) GetGroupMessages(ctx context.Context, req *pb.GetGroupMessagesRequest) (*pb.GetGroupMessagesResponse, error) {
	groupUUID, err := utils.ParseUUID(req.GroupId)
	if err != nil {
		return nil, err
	}

	userUUID, err := utils.ParseUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	isMember, err := u.itemRepo.IsGroupMember(ctx, groupUUID, userUUID)
	if err != nil {
		return nil, err
	}

	if !isMember {
		return nil, errors.New("unauthorized")
	}

	msgs, err := u.msgRepo.GetGroupMessages(ctx, req.GroupId, req.UserId)
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
