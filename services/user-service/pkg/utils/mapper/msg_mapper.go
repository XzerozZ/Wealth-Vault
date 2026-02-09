package mapper

import (
	"wealth-vault/user-service/internal/domain"
	pb "wealth-vault/user-service/pkg/pb/proto/user"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToGroupMessageProto(d domain.GroupMessage, userID string) *pb.MessageDetail {
	senderName := "Unknown"
	senderImg := ""
	if d.Sender != nil {
		senderName = d.Sender.Username
		senderImg = d.Sender.Profile
	}

	return &pb.MessageDetail{
		Id:          d.ID.String(),
		SenderId:    d.SenderID.String(),
		MsgType:     d.MsgType,
		Content:     d.Content,
		Metadata:    d.Metadata,
		CreatedAt:   timestamppb.New(d.CreatedAt),
		SenderName:  senderName,
		SenderImage: senderImg,
		IsMe:        d.SenderID.String() == userID,
	}
}

func ToPrivateMessageProto(d domain.PrivateMessage, userID string) *pb.MessageDetail {
	senderName := "Unknown"
	senderImg := ""
	if d.Sender != nil {
		senderName = d.Sender.Username
		senderImg = d.Sender.Profile
	}

	return &pb.MessageDetail{
		Id:          d.ID.String(),
		SenderId:    d.SenderID.String(),
		MsgType:     d.MsgType,
		Content:     d.Content,
		Metadata:    d.Metadata,
		CreatedAt:   timestamppb.New(d.CreatedAt),
		SenderName:  senderName,
		SenderImage: senderImg,
		IsMe:        d.SenderID.String() == userID,
	}
}
