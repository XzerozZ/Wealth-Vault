package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type LandRepository interface {
	CreateLand(ctx context.Context, item *domain.Land) error
	GetLand(ctx context.Context, uid uuid.UUID) ([]*domain.Land, error)
	GetLandByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Land, error)
	GetLandByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*domain.Land, error)
	UpdateLand(ctx context.Context, item *domain.Land, addBuildIDs []uuid.UUID, removeBuildIDs []uuid.UUID) (*domain.Land, error)
	DeleteLand(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
}
