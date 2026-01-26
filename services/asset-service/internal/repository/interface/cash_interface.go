package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type CashRepository interface {
	CreateCash(ctx context.Context, cash *domain.Cash) error
	GetCash(ctx context.Context, uid uuid.UUID) ([]*domain.Cash, error)
	GetCashByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*domain.Cash, error)
	UpdateCash(ctx context.Context, cash *domain.Cash) (*domain.Cash, error)
	DeleteCash(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
}
