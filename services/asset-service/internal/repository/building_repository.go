package repository

import (
	"context"
	"time"
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

func (r *BuildingRepository) GetBuildingByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Building, error) {
	var items []*domain.Building
	if err := r.db.WithContext(ctx).Preload("Location").Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *BuildingRepository) GetBatchBuildingByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Building, error) {
	var items []*domain.Building
	if err := r.db.WithContext(ctx).Unscoped().Preload("Location").Preload("Files").Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *BuildingRepository) GetBuildingByID(ctx context.Context, id uuid.UUID) (*domain.Building, error) {
	var item domain.Building
	if err := r.db.WithContext(ctx).Preload("Files").Preload("Location").Preload("Lands").Preload("Insurances").First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *BuildingRepository) UpdateBuilding(ctx context.Context, item *domain.Building, addLandIDs, removeLandIDs, addInsIDs, removeInsIDs []uuid.UUID) (*domain.Building, error) {
	return item, r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(item).Omit("Files").Updates(item).Error; err != nil {
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

		if len(addInsIDs) > 0 {
			var instoAdd []domain.Insurance
			for _, id := range addInsIDs {
				instoAdd = append(instoAdd, domain.Insurance{ID: id})
			}

			if err := tx.Model(item).Association("Insurances").Append(instoAdd); err != nil {
				return err
			}
		}

		if len(removeInsIDs) > 0 {
			var insToDelete []domain.Insurance
			for _, id := range removeInsIDs {
				insToDelete = append(insToDelete, domain.Insurance{ID: id})
			}

			if err := tx.Model(item).Association("Insurances").Delete(insToDelete); err != nil {
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

		if err := tx.Preload("Files").Preload("Location").Preload("Lands").Preload("Insurances").First(item, "id = ?", item.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *BuildingRepository) SoftDeleteBuilding(ctx context.Context, id, uid uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, uid).Select("Location").Delete(&domain.Building{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *BuildingRepository) GetExpiredBuilding(ctx context.Context, olderThan time.Time) ([]domain.Building, error) {
	var buildings []domain.Building
	if err := r.db.WithContext(ctx).Unscoped().Preload("Files").Where("deleted_at IS NOT NULL AND deleted_at < ?", olderThan).Find(&buildings).Error; err != nil {
		return nil, err
	}

	return buildings, nil
}

func (r *BuildingRepository) HardDeleteBuilding(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var building domain.Building

		if err := tx.Unscoped().First(&building, "id = ?", id).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Model(&building).Association("Lands").Clear(); err != nil {
			return err
		}

		if err := tx.Unscoped().Model(&building).Association("Insurances").Clear(); err != nil {
			return err
		}

		if err := tx.Unscoped().Where("entity_id = ? AND entity_type = ?", id, "building").Delete(&domain.FileAssociate{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Delete(&building).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Delete(&domain.Location{}, building.LocationID).Error; err != nil {
			return err
		}

		return nil
	})
}
