package repository

import (
	"context"
	"fmt"
	"wealth-vault/user-service/internal/domain"

	"github.com/google/uuid"
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

func (r *MsgRepository) UpdateGrantMessageStatus(ctx context.Context, groupID, ownerID, targetID uuid.UUID, newMetadata string) error {
	query := `
		UPDATE group_messages 
		SET metadata = ?, updated_at = NOW()
		WHERE group_id = ? 
		AND sender_id = ? 
		AND metadata->>'target_user_id' = ? 
		AND metadata->>'type' = 'GRANT_ACCESS_PROMPT'
		AND (metadata->>'is_completed')::boolean = false
	`

	result := r.db.WithContext(ctx).Exec(query,
		newMetadata,
		groupID,
		ownerID,
		targetID.String(),
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no pending grant message found")
	}

	return nil
}
