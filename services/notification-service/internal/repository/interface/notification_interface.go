package reepository

import (
	"context"
	"wealth-vault/notification-service/internal/domain"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, item *domain.Notification) error
	GetByReceiver(ctx context.Context, receiverID uuid.UUID) ([]domain.Notification, error)
	MarkAsRead(ctx context.Context, notificationID uuid.UUID, receiverID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, receiverID uuid.UUID) error

	UpdateNotificationMetadata(ctx context.Context, targetID, senderID uuid.UUID, notiType string, metaUpdates map[string]interface{}) error
}
