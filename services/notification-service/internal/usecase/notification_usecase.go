package usecase

import (
	"context"
	"encoding/json"
	"time"
	"wealth-vault/notification-service/internal/domain"
	line "wealth-vault/notification-service/internal/infra/line"
	dispatch "wealth-vault/notification-service/internal/infra/push_provider/interface"
	push_provider "wealth-vault/notification-service/internal/infra/push_provider/interface"
	socket "wealth-vault/notification-service/internal/infra/socket"
	repo "wealth-vault/notification-service/internal/repository/interface"
	pb "wealth-vault/notification-service/pkg/pb/proto/auth"
	m "wealth-vault/notification-service/pkg/utils/message"

	"github.com/google/uuid"
)

type NotificationUsecase struct {
	repo       repo.NotificationRepository
	drepo      repo.DeviceRepository
	hub        socket.ISocketHub
	dispatch   dispatch.Dispatcher
	lineClient line.LineClient
	authClient pb.AuthServiceClient
}

func NewNotificationUsecase(
	repo repo.NotificationRepository,
	drepo repo.DeviceRepository,
	hub socket.ISocketHub,
	dispatch dispatch.Dispatcher,
	lineClient line.LineClient,
	authClient pb.AuthServiceClient,
) *NotificationUsecase {
	return &NotificationUsecase{
		repo:       repo,
		drepo:      drepo,
		hub:        hub,
		dispatch:   dispatch,
		lineClient: lineClient,
		authClient: authClient,
	}
}

func (u *NotificationUsecase) GetHistory(ctx context.Context, uid uuid.UUID) ([]domain.Notification, error) {
	history, err := u.repo.GetByReceiver(ctx, uid)
	if err != nil {
		return nil, err
	}

	return history, nil
}

func (u *NotificationUsecase) MarkAsRead(ctx context.Context, notiID uuid.UUID, receiverID uuid.UUID) error {
	return u.repo.MarkAsRead(ctx, notiID, receiverID)
}

func (u *NotificationUsecase) MarkAllAsRead(ctx context.Context, receiverID uuid.UUID) error {
	return u.repo.MarkAllAsRead(ctx, receiverID)
}

func (u *NotificationUsecase) emitToUser(userID string, event string, payload interface{}) {
	u.hub.Emit(userID, domain.WSMessage{
		Type:    "DATA_UPDATE",
		Event:   event,
		Payload: payload,
	})
}

func (u *NotificationUsecase) emitToGroup(groupID string, event string, payload interface{}) {
	u.hub.BroadcastToGroup(groupID, domain.WSMessage{
		Type:    "DATA_UPDATE",
		Event:   event,
		Payload: payload,
	})
}

func (u *NotificationUsecase) notifyTarget(
	ctx context.Context,
	receiverID uuid.UUID,
	senderID *uuid.UUID,
	entityType string,
	entityID uuid.UUID,
	message string,
	occurredAt int64,
	metadata map[string]interface{},
) error {
	metaStr := "{}"
	if metadata != nil {
		if b, err := json.Marshal(metadata); err == nil {
			metaStr = string(b)
		}
	}

	noti := &domain.Notification{
		ID:         uuid.New(),
		EntityType: entityType,
		EntityID:   entityID,
		Receiver:   receiverID,
		SenderID:   senderID,
		Message:    message,
		Metadata:   metaStr,
		Channel:    "IN_APP",
		CreatedAt:  time.Unix(occurredAt, 0),
		IsRead:     false,
	}

	if err := u.repo.CreateNotification(ctx, noti); err != nil {
		return err
	}

	receiverStr := receiverID.String()

	if u.hub.IsOnline(receiverStr) {
		u.hub.Emit(receiverStr, domain.WSMessage{
			Type:    "NOTIFICATION",
			Payload: noti,
		})
	} else {
		go u.sendPush(ctx, receiverID, message, noti)
	}

	return nil
}

func (u *NotificationUsecase) sendPush(ctx context.Context, userID uuid.UUID, message string, noti *domain.Notification) {
	tokens, err := u.drepo.GetActiveTokens(ctx, userID)
	if err != nil || len(tokens) == 0 {
		return
	}

	u.dispatch.SendToUser(ctx, tokens, push_provider.PushPayload{
		Title: m.EntityTypeToTitle(noti.EntityType),
		Body:  message,
		Data: map[string]string{
			"notification_id": noti.ID.String(),
			"entity_type":     noti.EntityType,
			"entity_id":       noti.EntityID.String(),
		},
	})
}

func (u *NotificationUsecase) notifyMany(
	ctx context.Context,
	targetIDs []string,
	senderUUID *uuid.UUID,
	entityType string,
	entityID uuid.UUID,
	message string,
	occurredAt int64,
) error {

	for _, idStr := range targetIDs {
		receiverID, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}

		if senderUUID != nil && *senderUUID == receiverID {
			continue
		}

		if err := u.notifyTarget(ctx, receiverID, senderUUID, entityType, entityID, message, occurredAt, nil); err != nil {
			return err
		}
	}

	return nil
}
