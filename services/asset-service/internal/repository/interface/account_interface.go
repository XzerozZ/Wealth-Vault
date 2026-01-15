package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, acc *domain.Account) error
}
