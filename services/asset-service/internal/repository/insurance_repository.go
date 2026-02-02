package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InsuranceRepository struct {
	db *gorm.DB
}

func NewInsuranceRepository(db *gorm.DB) *InsuranceRepository {
	return &InsuranceRepository{db: db}
}

func (r *InsuranceRepository) CreateInsurance(ctx context.Context, policy *domain.Insurance) error {
	if err := r.db.WithContext(ctx).Create(policy).Error; err != nil {
		return err
	}
	return nil
}

func (r *InsuranceRepository) GetInsurance(ctx context.Context, uid uuid.UUID) ([]*domain.Insurance, error) {
	var items []*domain.Insurance
	if err := r.db.WithContext(ctx).Where("user_id = ?", uid).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *InsuranceRepository) GetInsuranceByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Insurance, error) {
	var items []*domain.Insurance
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *InsuranceRepository) GetInsuranceByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*domain.Insurance, error) {
	var item domain.Insurance
	if err := r.db.WithContext(ctx).Preload("Files").First(&item, "id = ? AND user_id = ?", id, uid).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *InsuranceRepository) UpdateInsurance(ctx context.Context, policy *domain.Insurance) (*domain.Insurance, error) {
	return policy, r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(policy).Error; err != nil {
			return err
		}

		if err := tx.Preload("Files").First(policy, "id = ?", policy.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *InsuranceRepository) DeleteInsurance(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("entity_id = ? AND entity_type = ?", id, "insurance").Delete(&domain.FileAssociate{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ? AND user_id = ?", id, uid).Delete(&domain.Insurance{}).Error; err != nil {
			return err
		}

		return nil
	})
}
