package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, item *domain.Account) error
	GetAccount(ctx context.Context, uid uuid.UUID) ([]*domain.Account, error)
	GetAccountByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*domain.Account, error)
	UpdateAccount(ctx context.Context, item *domain.Account) (*domain.Account, error)
	DeleteAccount(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
}
