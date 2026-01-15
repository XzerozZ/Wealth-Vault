package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"gorm.io/gorm"
)

type CashRepository struct {
	db *gorm.DB
}

func NewCashRepository(db *gorm.DB) *CashRepository {
	return &CashRepository{db: db}
}

func (r *CashRepository) CreateCash(ctx context.Context, cash *domain.Cash) error {
	if err := r.db.WithContext(ctx).Create(&cash).Error; err != nil {
		return err
	}
	return nil
}

func (r *CashRepository) GetCash(ctx context.Context, uid string) ([]domain.Cash, error) {
	var cash []domain.Cash
	if err := r.db.WithContext(ctx).Where("user_id = ?", uid).Find(&cash).Error; err != nil {
		return nil, err
	}

	return cash, nil
}

func (r *CashRepository) GetCashByID(ctx context.Context, id string, uid string) (*domain.Cash, error) {
	var cash *domain.Cash
	if err := r.db.WithContext(ctx).Preload("Files").Where("id = ? AND user_id = ?", id, uid).Find(&cash).Error; err != nil {
		return nil, err
	}

	return cash, nil
}

func (r *CashRepository) UpdateCash(ctx context.Context, cash *domain.Cash, mask []string) (*domain.Cash, error) {
	tx := r.db.WithContext(ctx).Model(cash).Where("id = ? AND user_id = ?", cash.ID, cash.UserID)
	if len(mask) > 0 {
		tx = tx.Select(mask)
	}

	if err := tx.Updates(cash).Error; err != nil {
		return nil, err
	}

	if err := r.db.Preload("Files").First(cash, "id = ?", cash.ID).Error; err != nil {
		return nil, err
	}

	return cash, nil
}

func (r *CashRepository) DeleteCash(ctx context.Context, id string, uid string) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, uid).Delete(&domain.Cash{}).Error; err != nil {
		return err
	}

	return nil
}
