package repository

import (
	"context"
	"time"
	"wealth-vault/notification-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) RegisterDevice(ctx context.Context, req *domain.DeviceToken) error {
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "token"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"user_id", "platform", "device_name", "is_active", "updated_at",
		}),
	}).Create(req).Error; err != nil {
		return err
	}

	return nil
}

func (r *DeviceRepository) GetActiveTokens(ctx context.Context, userID uuid.UUID) ([]domain.DeviceToken, error) {
	var tokens []domain.DeviceToken
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = true", userID).
		Find(&tokens).Error; err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r *DeviceRepository) UnregisterDevice(ctx context.Context, userID uuid.UUID, token string) error {
	if err := r.db.WithContext(ctx).Model(&domain.DeviceToken{}).Where("user_id = ? AND token = ?", userID, token).
		Updates(map[string]interface{}{"is_active": false, "updated_at": time.Now()}).Error; err != nil {
		return nil
	}

	return nil
}

func (r *DeviceRepository) MarkTokenInactive(ctx context.Context, token string) error {
	if err := r.db.WithContext(ctx).Model(&domain.DeviceToken{}).Where("token = ?", token).
		Updates(map[string]interface{}{"is_active": false, "updated_at": time.Now()}).Error; err != nil {
		return nil
	}

	return nil
}
