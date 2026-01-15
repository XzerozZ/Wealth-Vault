package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, acc *domain.Account) error {
	if err := r.db.WithContext(ctx).Create(&acc).Error; err != nil {
		return err
	}
	return nil
}
