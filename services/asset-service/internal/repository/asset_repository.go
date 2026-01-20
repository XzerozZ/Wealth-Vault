package repository

import (
	"context"
	"fmt"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) CreateAsset(ctx context.Context, asset *domain.Asset) error {
	fmt.Println(asset.IsIncludedInNetWorth)
	if err := r.db.WithContext(ctx).Create(&asset).Error; err != nil {
		return err
	}
	return nil
}

func (r *AssetRepository) GetAsset(ctx context.Context, uid uuid.UUID) ([]*domain.Asset, error) {
	var assets []*domain.Asset
	if err := r.db.WithContext(ctx).Where("user_id = ?", uid).Find(&assets).Error; err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *AssetRepository) GetAssetByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*domain.Asset, error) {
	var asset domain.Asset
	if err := r.db.WithContext(ctx).Preload("Files").First(&asset, "id = ? AND user_id = ?", id, uid).Error; err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *AssetRepository) UpdateAsset(ctx context.Context, asset *domain.Asset, mask []string) (*domain.Asset, error) {
	return asset, r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		query := tx.Model(asset).Where("id = ? AND user_id = ?", asset.ID, asset.UserID)
		if len(mask) > 0 {
			query = query.Select(mask)
		}

		if err := query.Updates(asset).Error; err != nil {
			return err
		}

		if err := tx.Preload("Files").First(asset, "id = ?", asset.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *AssetRepository) DeleteAsset(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("entity_id = ? AND entity_type = ?", id, "asset").Delete(&domain.FileAssociate{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ? AND user_id = ?", id, uid).Delete(&domain.Asset{}).Error; err != nil {
			return err
		}

		return nil
	})
}
