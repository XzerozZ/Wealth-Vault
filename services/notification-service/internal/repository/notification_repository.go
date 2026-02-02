package repository

import (
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

func (r *NotificationRepository) CreateNotification(item *domain.Notification) error {
	if err := r.db.Create(item).Error; err != nil {
		return err
	}
	return nil
}

func (r *NotificationRepository) GetByReceiver(receiverID uuid.UUID) ([]domain.Notification, error) {
	var list []domain.Notification
	err := r.db.Where("receiver = ?", receiverID).
		Order("created_at desc").
		Limit(50).
		Find(&list).Error

	return list, err
}
