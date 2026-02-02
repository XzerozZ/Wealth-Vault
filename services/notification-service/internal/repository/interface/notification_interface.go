package reepository

import (
	"wealth-vault/notification-service/internal/domain"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	CreateNotification(item *domain.Notification) error
	GetByReceiver(receiverID uuid.UUID) ([]domain.Notification, error)
}
