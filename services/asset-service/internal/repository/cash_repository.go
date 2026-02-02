package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CashRepository struct {
	db *gorm.DB
}

func NewCashRepository(db *gorm.DB) *CashRepository {
	return &CashRepository{db: db}
}

func (r *CashRepository) CreateCash(ctx context.Context, cash *domain.Cash) error {
	if err := r.db.WithContext(ctx).Create(cash).Error; err != nil {
		return err
	}
	return nil
}

func (r *CashRepository) GetCash(ctx context.Context, uid uuid.UUID) ([]*domain.Cash, error) {
	var items []*domain.Cash
	if err := r.db.WithContext(ctx).Where("user_id = ?", uid).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *CashRepository) GetCashByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Cash, error) {
	var items []*domain.Cash
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *CashRepository) GetCashByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*domain.Cash, error) {
	var item domain.Cash
	if err := r.db.WithContext(ctx).Preload("Files").First(&item, "id = ? AND user_id = ?", id, uid).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *CashRepository) UpdateCash(ctx context.Context, cash *domain.Cash) (*domain.Cash, error) {
	return cash, r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(cash).Error; err != nil {
			return err
		}

		if err := tx.Preload("Files").First(cash, "id = ?", cash.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *CashRepository) DeleteCash(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("entity_id = ? AND entity_type = ?", id, "cash").Delete(&domain.FileAssociate{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ? AND user_id = ?", id, uid).Delete(&domain.Cash{}).Error; err != nil {
			return err
		}

		return nil
	})
}
