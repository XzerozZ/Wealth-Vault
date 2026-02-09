package repository

import (
	"context"
	"wealth-vault/user-service/internal/domain"

	"gorm.io/gorm"
)

type MsgRepository struct {
	db *gorm.DB
}

func NewMsgRepository(db *gorm.DB) *MsgRepository {
	return &MsgRepository{db: db}
}

func (r *MsgRepository) CreateMessage(ctx context.Context, log []domain.GroupMessage) error {
	return r.db.WithContext(ctx).Create(&log).Error
}

func (r *MsgRepository) CreatePrivateMessage(ctx context.Context, log []domain.PrivateMessage) error {
	return r.db.WithContext(ctx).Create(&log).Error
}

func (r *MsgRepository) GetGroupMessages(ctx context.Context, groupID string) ([]domain.GroupMessage, error) {
	var msgs []domain.GroupMessage
	if err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Order("created_at DESC").
		Preload("Sender").Find(&msgs).Error; err != nil {
		return nil, err
	}

	return msgs, nil
}

func (r *MsgRepository) GetPrivateMessages(ctx context.Context, userID, friendID string) ([]domain.PrivateMessage, error) {
	var msgs []domain.PrivateMessage
	if err := r.db.WithContext(ctx).
		Where(
			r.db.Where("sender_id = ? AND receiver_id = ?", userID, friendID).
				Or("sender_id = ? AND receiver_id = ?", friendID, userID),
		).
		Order("created_at DESC").Preload("Sender").Find(&msgs).Error; err != nil {
		return nil, err
	}

	return msgs, nil
}
