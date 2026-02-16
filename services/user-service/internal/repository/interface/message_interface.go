package repository

import (
	"context"
	"wealth-vault/user-service/internal/domain"
)

type MsgRepository interface {
	// ------ Create Message ------
	CreateMessage(ctx context.Context, log []domain.GroupMessage) error
	CreatePrivateMessage(ctx context.Context, log []domain.PrivateMessage) error

	// ------ Get Message ------
	GetGroupMessages(ctx context.Context, groupID string) ([]domain.GroupMessage, error)
	GetPrivateMessages(ctx context.Context, userID, friendID string) ([]domain.PrivateMessage, error)
}
