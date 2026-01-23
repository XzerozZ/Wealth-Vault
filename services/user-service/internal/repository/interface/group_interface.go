package repository

import (
	"context"
	"wealth-vault/user-service/internal/domain"

	"github.com/google/uuid"
)

type GroupRepository interface {
	CreateGroup(ctx context.Context, group *domain.Group, initialMembers []string) error
	GetMember(ctx context.Context, id uuid.UUID) ([]*domain.User, int64, error)
	GetGroup(ctx context.Context, id uuid.UUID) (*domain.Group, int64, error)
	IsUserMember(ctx context.Context, id uuid.UUID, userID uuid.UUID) (bool, error)
	UpdateGroup(ctx context.Context, group *domain.Group, mask []string) (*domain.Group, int64, error)
}
