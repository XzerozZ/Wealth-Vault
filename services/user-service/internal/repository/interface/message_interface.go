package repository

import (
	"context"
	"wealth-vault/user-service/internal/domain"

	"github.com/google/uuid"
)

type MsgRepository interface {
	// ------ Create Message ------
	CreateMessage(ctx context.Context, log []domain.GroupMessage) error
	CreatePrivateMessage(ctx context.Context, log []domain.PrivateMessage) error

	// ------ Get Message ------
	GetGroupMessages(ctx context.Context, groupID string, userID string) ([]domain.GroupMessage, error)
	GetPrivateMessages(ctx context.Context, userID, friendID string) ([]domain.PrivateMessage, error)

	UpdateGrantMessageStatus(ctx context.Context, groupID, ownerID, targetID uuid.UUID, newMetadata string) error
	CloseAllGrantPromptsForTarget(ctx context.Context, groupID, targetID uuid.UUID) error
}
