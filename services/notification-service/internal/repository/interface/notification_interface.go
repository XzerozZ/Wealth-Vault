package reepository

import (
	"context"
	"wealth-vault/notification-service/internal/domain"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, item *domain.Notification) error
	GetByReceiver(ctx context.Context, receiverID uuid.UUID) ([]domain.Notification, error)
}
