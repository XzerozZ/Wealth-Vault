package repository

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type LandRepository interface {
	CreateLand(ctx context.Context, item *domain.Land) error
	GetLand(ctx context.Context, uid uuid.UUID) ([]*domain.Land, error)
	GetLandByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Land, error)
	GetBatchLandByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Land, error)
	GetLandByID(ctx context.Context, id uuid.UUID) (*domain.Land, error)
	GetLandByUserID(ctx context.Context, uid uuid.UUID) ([]*domain.Land, error)
	UpdateLand(ctx context.Context, item *domain.Land, addBuildIDs []uuid.UUID, removeBuildIDs []uuid.UUID) (*domain.Land, error)
	SoftDeleteLand(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
	GetExpiredLand(ctx context.Context, olderThan time.Time) ([]domain.Land, error)
	HardDeleteLand(ctx context.Context, id uuid.UUID) error
}
