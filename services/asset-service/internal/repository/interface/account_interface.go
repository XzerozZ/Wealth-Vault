package repository

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, item *domain.Account) error
	GetAccount(ctx context.Context, uid uuid.UUID) ([]*domain.Account, error)
	GetAccountByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Account, error)
	GetBatchAccountByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Account, error)
	GetAccountByID(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	UpdateAccount(ctx context.Context, item *domain.Account) (*domain.Account, error)
	SoftDeleteAccount(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
	GetExpiredAccounts(ctx context.Context, olderThan time.Time) ([]domain.Account, error)
	HardDeleteAccount(ctx context.Context, id uuid.UUID) error
}
