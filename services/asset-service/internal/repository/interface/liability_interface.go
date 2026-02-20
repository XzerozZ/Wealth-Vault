package repository

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type LiabilityRepository interface {
	CreateLiability(ctx context.Context, asset *domain.Liability) error
	GetLiability(ctx context.Context, uid uuid.UUID) ([]*domain.Liability, error)
	GetLiabilityByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Liability, error)
	GetBatchLiabilityByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Liability, error)
	GetLiabilityByID(ctx context.Context, id uuid.UUID) (*domain.Liability, error)
	UpdateLiability(ctx context.Context, lia *domain.Liability) (*domain.Liability, error)
	SoftDeleteLiability(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
	GetExpiredLiability(ctx context.Context, olderThan time.Time) ([]domain.Liability, error)
	HardDeleteLiability(ctx context.Context, id uuid.UUID) error
}
