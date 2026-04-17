package repository

import (
	"context"
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

func (r *MsgRepository) GetGroupMessages(ctx context.Context, groupID string, userID string) ([]domain.GroupMessage, error) {
	var msgs []domain.GroupMessage
	if err := r.db.WithContext(ctx).
		Table("group_messages").
		Joins("JOIN group_members ON group_members.group_id::uuid = group_messages.group_id::uuid").
		Where("group_messages.group_id = ?::uuid", groupID).
		Where("group_members.user_id = ?::uuid", userID).
		Where("group_messages.created_at >= group_members.joined_at").
		Where(
			r.db.Where("group_messages.sender_id = ?::uuid", userID).
				Or("group_messages.msg_type IN ?", []string{"ASSET_CARD", "SYSTEM_ALERT"}),
		).
		Order("group_messages.created_at DESC").
		Preload("Sender").
		Find(&msgs).Error; err != nil {
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
        WITH target_msg AS (
            SELECT id FROM group_messages
            WHERE group_id = ? 
            AND sender_id = ? 
            AND metadata->>'target_user_id' = ? 
            AND metadata->>'type' = 'GRANT_ACCESS_PROMPT'
            AND (metadata->>'is_completed')::boolean = false
        )
        UPDATE group_messages
        SET metadata = ?, updated_at = NOW()
        FROM target_msg
        WHERE group_messages.id = target_msg.id
    `
	result := r.db.WithContext(ctx).Exec(query, groupID, ownerID, targetID.String(), newMetadata)
	return result.Error
}

func (r *MsgRepository) CloseAllGrantPromptsForTarget(ctx context.Context, groupID, targetID uuid.UUID) error {
	statusMeta := `{"is_action_required": true, "is_completed": true, "status": "member_left"}`

	query := `
        UPDATE group_messages 
        SET metadata = ?::jsonb,
            updated_at = NOW()
        WHERE group_id = ? 
        AND metadata->>'target_user_id' = ? 
        AND metadata->>'type' = 'GRANT_ACCESS_PROMPT'
        AND (metadata->>'is_completed')::boolean = false
    `

	return r.db.WithContext(ctx).Exec(query, statusMeta, groupID, targetID.String()).Error
}
