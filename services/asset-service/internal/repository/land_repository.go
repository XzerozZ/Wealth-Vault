package repository

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LandRepository struct {
	db *gorm.DB
}

func NewLandRepository(db *gorm.DB) *LandRepository {
	return &LandRepository{db: db}
}

func (r *LandRepository) CreateLand(ctx context.Context, item *domain.Land) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return err
	}

	return nil
}

func (r *LandRepository) GetLand(ctx context.Context, uid uuid.UUID) ([]*domain.Land, error) {
	var items []*domain.Land
	if err := r.db.WithContext(ctx).Preload("Location").Where("user_id = ?", uid).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *LandRepository) GetLandByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Land, error) {
	var items []*domain.Land
	if err := r.db.WithContext(ctx).Preload("Location").Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *LandRepository) GetBatchLandByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Land, error) {
	var items []*domain.Land
	if err := r.db.WithContext(ctx).Unscoped().Preload("Location").Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *LandRepository) GetLandByID(ctx context.Context, id uuid.UUID) (*domain.Land, error) {
	var item domain.Land
	if err := r.db.WithContext(ctx).Preload("Files").Preload("Location").Preload("Buildings").First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *LandRepository) UpdateLand(ctx context.Context, item *domain.Land, addBuildIDs []uuid.UUID, removeBuildIDs []uuid.UUID) (*domain.Land, error) {
	return item, r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(item).Omit("Files").Updates(item).Error; err != nil {
			return err
		}

		if err := tx.Save(&item.Location).Error; err != nil {
			return err
		}

		if len(removeBuildIDs) > 0 {
			var buildToDelete []domain.Building
			for _, id := range removeBuildIDs {
				buildToDelete = append(buildToDelete, domain.Building{ID: id})
			}

			if err := tx.Model(item).Association("Lands").Delete(buildToDelete); err != nil {
				return err
			}
		}

		if len(addBuildIDs) > 0 {
			var buildToAdd []domain.Building
			for _, id := range addBuildIDs {
				buildToAdd = append(buildToAdd, domain.Building{ID: id})
			}

			if err := tx.Model(item).Association("Buildings").Append(buildToAdd); err != nil {
				return err
			}
		}

		if err := tx.Preload("Files").Preload("Location").Preload("Buildings").First(item, "id = ?", item.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *LandRepository) SoftDeleteLand(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, uid).Delete(&domain.Land{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *LandRepository) GetExpiredLand(ctx context.Context, olderThan time.Time) ([]domain.Land, error) {
	var lands []domain.Land
	if err := r.db.WithContext(ctx).Unscoped().Preload("Files").Where("deleted_at IS NOT NULL AND deleted_at < ?", olderThan).Find(&lands).Error; err != nil {
		return nil, err
	}

	return lands, nil
}

func (r *LandRepository) HardDeleteLand(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var land domain.Land

		if err := tx.Unscoped().Where("id = ?", id).First(&land).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Model(&land).Association("Buildings").Clear(); err != nil {
			return err
		}

		if err := tx.Unscoped().Where("entity_id = ? AND entity_type = ?", id, "land").Delete(&domain.FileAssociate{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Select("Location").Delete(&land).Error; err != nil {
			return err
		}

		return nil
	})
}
