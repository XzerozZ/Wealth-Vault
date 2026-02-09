package usecase

import (
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
