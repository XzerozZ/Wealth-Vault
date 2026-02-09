package repository

import (
	"context"
	"time"
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

func (r *InsuranceRepository) GetBatchInsuranceByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Insurance, error) {
	var items []*domain.Insurance
	if err := r.db.WithContext(ctx).Unscoped().Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *InsuranceRepository) GetInsuranceByID(ctx context.Context, id uuid.UUID) (*domain.Insurance, error) {
	var item domain.Insurance
	if err := r.db.WithContext(ctx).Preload("Files").First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *InsuranceRepository) GetInsuranceByUserID(ctx context.Context, uid uuid.UUID) ([]*domain.Insurance, error) {
	var items []*domain.Insurance
	if err := r.db.WithContext(ctx).Where("user_id = ?", uid).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
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

func (r *InsuranceRepository) SoftDeleteInsurances(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, uid).Delete(&domain.Insurance{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *InsuranceRepository) GetExpiredInsurances(ctx context.Context, olderThan time.Time) ([]domain.Insurance, error) {
	var insurances []domain.Insurance
	if err := r.db.WithContext(ctx).Unscoped().Preload("Files").Where("deleted_at IS NOT NULL AND deleted_at < ?", olderThan).Find(&insurances).Error; err != nil {
		return nil, err
	}

	return insurances, nil
}

func (r *InsuranceRepository) HardDeleteInsurances(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var insurance domain.Insurance
		if err := tx.Unscoped().First(&insurance, "id = ?", id).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Model(&insurance).Association("Buildings").Clear(); err != nil {
			return err
		}

		if err := tx.Unscoped().Where("entity_id = ? AND entity_type = ?", id, "insurance").Delete(&domain.FileAssociate{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Delete(&insurance).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *InsuranceRepository) GetExpiringInsurances(ctx context.Context, days int) ([]*domain.Insurance, error) {
	var insurances []*domain.Insurance
	targetDate := time.Now().AddDate(0, 0, days).Format("2006-01-02")

	if err := r.db.WithContext(ctx).Where("DATE(exp_date) = ?", targetDate).Find(&insurances).Error; err != nil {
		return nil, err
	}

	return insurances, nil
}
