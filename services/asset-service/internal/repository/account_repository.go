package repository

import (
	"context"
	"time"
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

func (r *AccountRepository) GetAccountByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Account, error) {
	var items []*domain.Account
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *AccountRepository) GetBatchAccountByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Account, error) {
	var items []*domain.Account
	if err := r.db.WithContext(ctx).Unscoped().Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *AccountRepository) GetAccountByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	var item domain.Account
	if err := r.db.WithContext(ctx).Preload("Files").First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *AccountRepository) UpdateAccount(ctx context.Context, item *domain.Account) (*domain.Account, error) {
	return item, r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(item).Omit("Files").Updates(item).Error; err != nil {
			return err
		}

		if err := tx.Preload("Files").First(item, "id = ?", item.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *AccountRepository) SoftDeleteAccount(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, uid).Delete(&domain.Account{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *AccountRepository) GetExpiredAccounts(ctx context.Context, olderThan time.Time) ([]domain.Account, error) {
	var accounts []domain.Account
	if err := r.db.WithContext(ctx).Unscoped().Preload("Files").Where("deleted_at IS NOT NULL AND deleted_at < ?", olderThan).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *AccountRepository) HardDeleteAccount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("entity_id = ? AND entity_type = ?", id, "account").Delete(&domain.FileAssociate{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Where("id = ?", id).Delete(&domain.Account{}).Error; err != nil {
			return err
		}

		return nil
	})
}
