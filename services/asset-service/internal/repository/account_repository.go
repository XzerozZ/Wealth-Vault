package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, item *domain.Account) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) GetAccount(ctx context.Context, uid uuid.UUID) ([]*domain.Account, error) {
	var items []*domain.Account
	if err := r.db.WithContext(ctx).Where("user_id = ?", uid).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *AccountRepository) GetAccountByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*domain.Account, error) {
	var item domain.Account
	if err := r.db.WithContext(ctx).Preload("Files").First(&item, "id = ? AND user_id = ?", id, uid).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *AccountRepository) UpdateAccount(ctx context.Context, item *domain.Account) (*domain.Account, error) {
	return item, r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(item).Error; err != nil {
			return err
		}

		if err := tx.Preload("Files").First(item, "id = ?", item.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *AccountRepository) DeleteAccount(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("entity_id = ? AND entity_type = ?", id, "account").Delete(&domain.FileAssociate{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ? AND user_id = ?", id, uid).Delete(&domain.Account{}).Error; err != nil {
			return err
		}

		return nil
	})
}
