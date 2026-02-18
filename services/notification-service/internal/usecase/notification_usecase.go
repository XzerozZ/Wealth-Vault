package usecase

import (
	"context"
	"wealth-vault/notification-service/internal/domain"
	socket "wealth-vault/notification-service/internal/infra/socket/interface"
	repo "wealth-vault/notification-service/internal/repository/interface"

	"github.com/google/uuid"
)

type NotificationUsecase struct {
	repo repo.NotificationRepository
	hub  socket.ISocketHub
}

func NewNotificationUsecase(repo repo.NotificationRepository, hub socket.ISocketHub) *NotificationUsecase {
	return &NotificationUsecase{
		repo: repo,
		hub:  hub,
	}
}

func (u *NotificationUsecase) GetHistory(ctx context.Context, uid uuid.UUID) ([]domain.Notification, error) {
	history, err := u.repo.GetByReceiver(ctx, uid)
	if err != nil {
		return nil, err
	}

	return history, nil
}
