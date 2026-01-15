package usecase

import (
	"context"
	"wealth-vault/asset-service/internal/domain"
)

type CashUsecase interface {
	CreateCash(ctx context.Context, cash *domain.Cash) (string, error)
	GetCash(ctx context.Context, uid string) ([]domain.Cash, error)
	GetCashByID(ctx context.Context, id string, uid string) (*domain.Cash, error)
	DeleteCash(ctx context.Context, id string, uid string) error
}
