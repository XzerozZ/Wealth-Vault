package repository

import (
	"context"
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

func (r *LiabilityRepository) GetLiabilityByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*domain.Liability, error) {
	var lia domain.Liability
	if err := r.db.WithContext(ctx).Preload("Files").First(&lia, "id = ? AND user_id = ?", id, uid).Error; err != nil {
		return nil, err
	}

	return &lia, nil
}

func (r *LiabilityRepository) UpdateLiability(ctx context.Context, lia *domain.Liability) (*domain.Liability, error) {
	return lia, r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(lia).Error; err != nil {
			return err
		}

		if err := tx.Preload("Files").First(lia, "id = ?", lia.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *LiabilityRepository) DeleteLiability(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("entity_id = ? AND entity_type = ?", id, "liability").Delete(&domain.FileAssociate{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ? AND user_id = ?", id, uid).Delete(&domain.Liability{}).Error; err != nil {
			return err
		}

		return nil
	})
}
