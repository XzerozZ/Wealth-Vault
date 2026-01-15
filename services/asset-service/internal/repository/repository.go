package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"gorm.io/gorm"
)

type AssetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) CreateCash(ctx context.Context, cash *domain.Cash) error {
	if err := r.db.WithContext(ctx).Create(&cash).Error; err != nil {
		return err
	}
	return nil
}

func (r *AssetRepository) GetCash(ctx context.Context, uid string) ([]domain.Cash, error) {
	var cash []domain.Cash
	if err := r.db.WithContext(ctx).Where("user_id = ?", uid).Find(&cash).Error; err != nil {
		return nil, err
	}

	return cash, nil
}

func (r *AssetRepository) GetCashByID(ctx context.Context, id string, uid string) (*domain.Cash, error) {
	var cash *domain.Cash
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, uid).Find(&cash).Error; err != nil {
		return nil, err
	}

	return cash, nil
}

func (r *AssetRepository) UpdateCash(ctx context.Context, cash *domain.Cash, mask []string) (*domain.Cash, error) {
	tx := r.db.WithContext(ctx).Model(cash).Where("id = ? AND user_id = ?", cash.ID, cash.UserID)
	if len(mask) > 0 {
		tx = tx.Select(mask)
	}

	if err := tx.Updates(cash).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(cash, "id = ?", cash.ID).Error; err != nil {
		return nil, err
	}

	return cash, nil
}
