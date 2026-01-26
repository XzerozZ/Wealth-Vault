package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type InvestmentRepository interface {
	CreateInvestment(ctx context.Context, invest *domain.Investment) error
	GetInvestment(ctx context.Context, uid uuid.UUID) ([]*domain.Investment, error)
	GetInvestmentByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*domain.Investment, error)
	UpdateInvestment(ctx context.Context, invest *domain.Investment) (*domain.Investment, error)
	DeleteInvestment(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
}
