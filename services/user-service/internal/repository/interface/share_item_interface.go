package repository

import (
	"context"
	"wealth-vault/user-service/internal/domain"

	"github.com/google/uuid"
)

type ShareItemRepository interface {
	ShareItemtoGroup(ctx context.Context, items []domain.GroupItem) error
	ShareItemtoFriend(ctx context.Context, items []domain.FriendItem) error
	ShareItemtoEmail(ctx context.Context, items []domain.EmailItem) error
	GetExistingSharedMap(ctx context.Context, ownerID, friendID uuid.UUID) (map[string]bool, error)
	IsItemSharedtoGroup(ctx context.Context, groupID, entityID uuid.UUID, entityType string) (bool, error)
	IsItemSharedtoFriend(ctx context.Context, friendID, entityID uuid.UUID, entityType string) (bool, error)
	IsItemSharedtoEmail(ctx context.Context, entityID uuid.UUID, email, entityType string) (bool, error)
	GetSharedIteminGroup(ctx context.Context, groupID, userID uuid.UUID) ([]domain.GroupItem, error)
	GetSharedIteminFriend(ctx context.Context, friendID, userID uuid.UUID) ([]domain.FriendItem, error)
	DeleteIteminGroup(ctx context.Context, itemID uuid.UUID, userID uuid.UUID) error
	DeleteIteminFriend(ctx context.Context, itemID uuid.UUID, userID uuid.UUID) error
	GetPendingEmails(ctx context.Context) ([]domain.EmailItem, error)
	MarkEmailsAsSent(ctx context.Context, ids []uuid.UUID) error
	IsGroupMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) (bool, error)
	GetOwnedItemIDs(ctx context.Context, itemIDs []string, ownerID uuid.UUID) ([]uuid.UUID, error)
	CountItemsByOwner(ctx context.Context, itemIDs []string, ownerID uuid.UUID) (int64, error)
	AddMember(ctx context.Context, members []domain.GroupMember) error
	BatchCreateViewers(ctx context.Context, viewers []domain.GroupItemViewer) error
	GetFutureItemsInGroup(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error)
	GetItemOwnersInGroup(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error)
	DeleteAllReferencesByEntityID(ctx context.Context, entityID uuid.UUID) error
}
