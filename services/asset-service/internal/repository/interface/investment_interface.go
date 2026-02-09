package repository

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type InvestmentRepository interface {
	CreateInvestment(ctx context.Context, invest *domain.Investment) error
	GetInvestment(ctx context.Context, uid uuid.UUID) ([]*domain.Investment, error)
	GetInvestmentByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Investment, error)
	GetBatchInvestmentByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Investment, error)
	GetInvestmentByID(ctx context.Context, id uuid.UUID) (*domain.Investment, error)
	GetInvestmentByUserID(ctx context.Context, uid uuid.UUID) ([]*domain.Investment, error)
	UpdateInvestment(ctx context.Context, invest *domain.Investment) (*domain.Investment, error)
	SoftDeleteInvestment(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
	GetExpiredInvestment(ctx context.Context, olderThan time.Time) ([]domain.Investment, error)
	HardDeleteInvestment(ctx context.Context, id uuid.UUID) error
}
