package repository

import (
	"context"
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

func (r *AssetRepository) CheckExists(ctx context.Context, entityType string, id uuid.UUID, uid uuid.UUID) (bool, error) {
	var count int64
	var err error

	switch entityType {
	case "account":
		err = r.db.WithContext(ctx).Model(&domain.Account{}).Where("id = ? AND user_id = ?", id, uid).Count(&count).Error
	case "investment":
		err = r.db.WithContext(ctx).Model(&domain.Investment{}).Where("id = ? AND user_id = ?", id, uid).Count(&count).Error
	case "insurance":
		err = r.db.WithContext(ctx).Model(&domain.Insurance{}).Where("id = ? AND user_id = ?", id, uid).Count(&count).Error
	case "building":
		err = r.db.WithContext(ctx).Model(&domain.Building{}).Where("id = ? AND user_id = ?", id, uid).Count(&count).Error
	case "land":
		err = r.db.WithContext(ctx).Model(&domain.Land{}).Where("id = ? AND user_id = ?", id, uid).Count(&count).Error
	case "cash":
		err = r.db.WithContext(ctx).Model(&domain.Cash{}).Where("id = ? AND user_id = ?", id, uid).Count(&count).Error
	case "liability":
		err = r.db.WithContext(ctx).Model(&domain.Liability{}).Where("id = ? AND user_id = ?", id, uid).Count(&count).Error
	default:
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
