package repository

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type CashRepository interface {
	CreateCash(ctx context.Context, cash *domain.Cash) error
	GetCash(ctx context.Context, uid uuid.UUID) ([]*domain.Cash, error)
	GetCashByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Cash, error)
	GetBatchCashByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Cash, error)
	GetCashByID(ctx context.Context, id uuid.UUID) (*domain.Cash, error)
	GetCashByUserID(ctx context.Context, uid uuid.UUID) ([]*domain.Cash, error)
	UpdateCash(ctx context.Context, cash *domain.Cash) (*domain.Cash, error)
	SoftDeleteCash(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
	HardDeleteCash(ctx context.Context, id uuid.UUID) error
	GetExpiredCash(ctx context.Context, olderThan time.Time) ([]domain.Cash, error)
}
