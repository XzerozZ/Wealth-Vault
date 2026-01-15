package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"
)

type AssetRepository interface {
	CreateCash(ctx context.Context, cash *domain.Cash) error
	GetCash(ctx context.Context, uid string) ([]domain.Cash, error)
	GetCashByID(ctx context.Context, id string, uid string) (*domain.Cash, error)
	UpdateCash(ctx context.Context, cash *domain.Cash, mask []string) (*domain.Cash, error)
}
