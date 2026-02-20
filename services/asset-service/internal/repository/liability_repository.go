package repository

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LiabilityRepository struct {
	db *gorm.DB
}

func NewLiabilityRepository(db *gorm.DB) *LiabilityRepository {
	return &LiabilityRepository{db: db}
}

func (r *LiabilityRepository) CreateLiability(ctx context.Context, lia *domain.Liability) error {
	if err := r.db.WithContext(ctx).Create(&lia).Error; err != nil {
		return err
	}
	return nil
}

func (r *LiabilityRepository) GetLiability(ctx context.Context, uid uuid.UUID) ([]*domain.Liability, error) {
	var lias []*domain.Liability
	if err := r.db.WithContext(ctx).Where("user_id = ?", uid).Find(&lias).Error; err != nil {
		return nil, err
	}

	return lias, nil
}

func (r *LiabilityRepository) GetLiabilityByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Liability, error) {
	var items []*domain.Liability
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *LiabilityRepository) GetBatchLiabilityByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Liability, error) {
	var items []*domain.Liability
	if err := r.db.WithContext(ctx).Unscoped().Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *LiabilityRepository) GetLiabilityByID(ctx context.Context, id uuid.UUID) (*domain.Liability, error) {
	var lia domain.Liability
	if err := r.db.WithContext(ctx).Preload("Files").First(&lia, "id = ? ", id).Error; err != nil {
		return nil, err
	}

	return &lia, nil
}

func (r *LiabilityRepository) UpdateLiability(ctx context.Context, lia *domain.Liability) (*domain.Liability, error) {
	return lia, r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(lia).Omit("Files").Updates(lia).Error; err != nil {
			return err
		}

		if err := tx.Preload("Files").First(lia, "id = ?", lia.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *LiabilityRepository) SoftDeleteLiability(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, uid).Delete(&domain.Liability{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *LiabilityRepository) GetExpiredLiability(ctx context.Context, olderThan time.Time) ([]domain.Liability, error) {
	var lias []domain.Liability
	if err := r.db.WithContext(ctx).Unscoped().Preload("Files").Where("deleted_at IS NOT NULL AND deleted_at < ?", olderThan).Find(&lias).Error; err != nil {
		return nil, err
	}

	return lias, nil
}

func (r *LiabilityRepository) HardDeleteLiability(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("entity_id = ? AND entity_type = ?", id, "liability").Delete(&domain.FileAssociate{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Where("id = ?", id).Delete(&domain.Liability{}).Error; err != nil {
			return err
		}

		return nil
	})
}
