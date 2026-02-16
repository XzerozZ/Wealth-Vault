package repository

import (
	"context"
	"wealth-vault/user-service/internal/domain"

	"github.com/google/uuid"
)

type GroupRepository interface {
	// ------ Create Group and Group log ------
	CreateGroup(ctx context.Context, group *domain.Group, initialMembers []string) error
	CreateLog(ctx context.Context, log *domain.GroupLog) error

	// ------ Get Member and Group ------
	GetMember(ctx context.Context, id uuid.UUID) ([]*domain.User, int64, error)
	GetGroup(ctx context.Context, id uuid.UUID) (*domain.Group, int64, error)
	AllGetGroup(ctx context.Context, uid uuid.UUID) ([]domain.GroupWithCount, error)

	// ------ Check Member in Group ------
	IsUserMember(ctx context.Context, id uuid.UUID, userID uuid.UUID) (bool, error)

	// ------ Update Group ------
	UpdateGroup(ctx context.Context, group *domain.Group, mask []string, logEntry *domain.GroupLog) (*domain.Group, int64, error)

	// ------ Delete and Remove ------
	RemoveMemberAndTheirSharedItems(ctx context.Context, groupID, memberID uuid.UUID, logEntry *domain.GroupLog) error
	DeleteGroup(ctx context.Context, groupID uuid.UUID) error
}
