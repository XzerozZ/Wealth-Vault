package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BuildingRepository struct {
	db *gorm.DB
}

func NewBuildingRepository(db *gorm.DB) *BuildingRepository {
	return &BuildingRepository{db: db}
}

func (r *BuildingRepository) CreateBuilding(ctx context.Context, item *domain.Building) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return err
	}

	return nil
}

func (r *BuildingRepository) GetBuilding(ctx context.Context, uid uuid.UUID) ([]*domain.Building, error) {
	var items []*domain.Building
	if err := r.db.WithContext(ctx).Preload("Location").Where("user_id = ?", uid).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *BuildingRepository) GetBuildingByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*domain.Building, error) {
	var item domain.Building
	if err := r.db.WithContext(ctx).Preload("Files").Preload("Location").Preload("Lands").First(&item, "id = ? AND user_id = ?", id, uid).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *BuildingRepository) UpdateBuilding(ctx context.Context, item *domain.Building, addLandIDs []uuid.UUID, removeLandIDs []uuid.UUID) (*domain.Building, error) {
	return item, r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(item).Error; err != nil {
			return err
		}

		if err := tx.Save(&item.Location).Error; err != nil {
			return err
		}

		if len(removeLandIDs) > 0 {
			var landsToDelete []domain.Land
			for _, id := range removeLandIDs {
				landsToDelete = append(landsToDelete, domain.Land{ID: id})
			}

			if err := tx.Model(item).Association("Lands").Delete(landsToDelete); err != nil {
				return err
			}
		}

		if len(addLandIDs) > 0 {
			var landsToAdd []domain.Land
			for _, id := range addLandIDs {
				landsToAdd = append(landsToAdd, domain.Land{ID: id})
			}

			if err := tx.Model(item).Association("Lands").Append(landsToAdd); err != nil {
				return err
			}
		}

		if err := tx.Preload("Files").Preload("Location").Preload("Lands").First(item, "id = ?", item.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *BuildingRepository) DeleteBuilding(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var building domain.Building
		if err := tx.Where("id = ? AND user_id = ?", id, uid).First(&building).Error; err != nil {
			return err
		}

		if err := tx.Model(&building).Association("Lands").Clear(); err != nil {
			return err
		}

		if err := tx.Where("entity_id = ? AND entity_type = ?", id, "building").Delete(&domain.FileAssociate{}).Error; err != nil {
			return err
		}

		if err := tx.Select("Location").Delete(&building).Error; err != nil {
			return err
		}

		return nil
	})
}
