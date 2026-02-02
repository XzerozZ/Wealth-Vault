package usecase

import (
	"fmt"
	"log"
	"time"
	"wealth-vault/notification-service/internal/domain"
	"wealth-vault/notification-service/internal/infra/socket"
	repo "wealth-vault/notification-service/internal/repository"

	"github.com/google/uuid"
)

type NotificationUsecase struct {
	repo *repo.NotificationRepository
	hub  *socket.SocketHub
}

func NewNotificationUsecase(repo *repo.NotificationRepository, hub *socket.SocketHub) *NotificationUsecase {
	return &NotificationUsecase{
		repo: repo,
		hub:  hub,
	}
}

func (u *NotificationUsecase) HandleGroupMemberAdded(evt domain.GroupMemberAddedEvent) {
	groupID, _ := uuid.Parse(evt.GroupID)
	occurredTime := time.Unix(evt.OccurredAt, 0)

	var senderUUIDPtr *uuid.UUID
	if evt.SenderID != "" {
		if id, err := uuid.Parse(evt.SenderID); err == nil {
			senderUUIDPtr = &id
		}
	}

	for _, targetIDStr := range evt.TargetUserIDs {
		receiverID, err := uuid.Parse(targetIDStr)
		if err != nil {
			continue
		}

		if senderUUIDPtr != nil && *senderUUIDPtr == receiverID {
			continue
		}

		noti := &domain.Notification{
			ID:         uuid.New(),
			EntityType: "GROUP",
			EntityID:   groupID,
			Receiver:   receiverID,
			SenderID:   senderUUIDPtr,
			Channel:    "IN_APP",
			Message:    "มีสมาชิกใหม่ได้ถูกเชิญเข้ามาในกลุ่ม มาให้สิทธิ์ในการมองเห็นทรัพย์สินให้กันเถอะ 💸",
			CreatedAt:  occurredTime,
		}

		if err := u.repo.CreateNotification(noti); err != nil {
			log.Printf("❌ Save DB Error: %v", err)
		}

		u.hub.Emit(targetIDStr, noti)
	}
}

func (u *NotificationUsecase) HandleItemShared(evt domain.ItemSharedEvent) {
	occurredTime := time.Unix(evt.OccurredAt, 0)
	assetUUID, _ := uuid.Parse(evt.AssetID)

	var senderUUIDPtr *uuid.UUID
	if evt.SenderID != "" {
		if id, err := uuid.Parse(evt.SenderID); err == nil {
			senderUUIDPtr = &id
		}
	}

	for _, targetIDStr := range evt.TargetUserIDs {
		receiverID, err := uuid.Parse(targetIDStr)
		if err != nil {
			continue
		}

		if senderUUIDPtr != nil && *senderUUIDPtr == receiverID {
			continue
		}

		noti := &domain.Notification{
			ID:         uuid.New(),
			EntityType: "ASSET",
			EntityID:   assetUUID,
			Receiver:   receiverID,
			SenderID:   senderUUIDPtr,
			Channel:    "IN_APP",
			Message:    fmt.Sprintf("%s ได้แชร์รายการใหม่ให้กับคุณ", evt.SenderName),
			CreatedAt:  occurredTime,
		}

		if err := u.repo.CreateNotification(noti); err != nil {
			log.Printf("❌ Save DB Error: %v", err)
		}

		u.hub.Emit(targetIDStr, noti)
	}
}

func (u *NotificationUsecase) GetHistory(userIDStr string) ([]domain.Notification, error) {
	uid, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	history, err := u.repo.GetByReceiver(uid)
	if err != nil {
		return nil, err
	}

	return history, nil
}
