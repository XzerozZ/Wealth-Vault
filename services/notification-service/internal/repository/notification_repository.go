package repository

import (
	"context"
	"encoding/json"
	"wealth-vault/notification-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) CreateNotification(ctx context.Context, item *domain.Notification) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) GetByReceiver(ctx context.Context, receiverID uuid.UUID) ([]domain.Notification, error) {
	var list []domain.Notification
	if err := r.db.WithContext(ctx).Where("receiver = ?", receiverID).
		Order("created_at desc").Limit(50).Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID uuid.UUID, receiverID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&domain.Notification{}).Where("id = ? AND receiver = ?", notificationID, receiverID).
		Update("is_read", true).Error; err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, receiverID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&domain.Notification{}).Where("receiver = ? AND is_read = false", receiverID).
		Update("is_read", true).Error; err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) UpdateNotificationMetadata(ctx context.Context, targetID, senderID uuid.UUID, notiType string, metaUpdates map[string]interface{}) error {
	metaJSON, _ := json.Marshal(metaUpdates)
	query := `
		UPDATE notifications 
		SET metadata = metadata || ?::jsonb,
			updated_at = NOW()
		WHERE receiver = ? 
		AND sender_id = ? 
		AND entity_type = ?
		AND (metadata->>'is_completed')::boolean = false
	`

	return r.db.WithContext(ctx).Exec(query, string(metaJSON), targetID, senderID, notiType).Error
}
