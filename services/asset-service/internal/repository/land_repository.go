package repository

import (
	"context"
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

func (r *LandRepository) GetLandByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*domain.Land, error) {
	var item domain.Land
	if err := r.db.WithContext(ctx).Preload("Files").Preload("Location").Preload("Buildings").First(&item, "id = ? AND user_id = ?", id, uid).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *LandRepository) UpdateLand(ctx context.Context, item *domain.Land, addBuildIDs []uuid.UUID, removeBuildIDs []uuid.UUID) (*domain.Land, error) {
	return item, r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(item).Error; err != nil {
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

func (r *LandRepository) DeleteLand(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var land domain.Land
		if err := tx.Where("id = ? AND user_id = ?", id, uid).First(&land).Error; err != nil {
			return err
		}

		if err := tx.Model(&land).Association("Buildings").Clear(); err != nil {
			return err
		}

		if err := tx.Where("entity_id = ? AND entity_type = ?", id, "land").Delete(&domain.FileAssociate{}).Error; err != nil {
			return err
		}

		if err := tx.Select("Location").Delete(&land).Error; err != nil {
			return err
		}

		return nil
	})
}
