package repository

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvestmentRepository struct {
	db *gorm.DB
}

func NewInvestmentRepository(db *gorm.DB) *InvestmentRepository {
	return &InvestmentRepository{db: db}
}

func (r *InvestmentRepository) CreateInvestment(ctx context.Context, invest *domain.Investment) error {
	if err := r.db.WithContext(ctx).Create(invest).Error; err != nil {
		return err
	}
	return nil
}

func (r *InvestmentRepository) GetInvestment(ctx context.Context, uid uuid.UUID) ([]*domain.Investment, error) {
	var items []*domain.Investment
	if err := r.db.WithContext(ctx).Where("user_id = ?", uid).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *InvestmentRepository) GetInvestmentByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Investment, error) {
	var items []*domain.Investment
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *InvestmentRepository) GetBatchInvestmentByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Investment, error) {
	var items []*domain.Investment
	if err := r.db.WithContext(ctx).Unscoped().Where("id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *InvestmentRepository) GetInvestmentByID(ctx context.Context, id uuid.UUID) (*domain.Investment, error) {
	var item domain.Investment
	if err := r.db.WithContext(ctx).Preload("Files").First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *InvestmentRepository) GetInvestmentByUserID(ctx context.Context, uid uuid.UUID) ([]*domain.Investment, error) {
	var items []*domain.Investment
	if err := r.db.WithContext(ctx).Where("user_id = ?", uid).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *InvestmentRepository) UpdateInvestment(ctx context.Context, invest *domain.Investment) (*domain.Investment, error) {
	return invest, r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(invest).Error; err != nil {
			return err
		}

		if err := tx.Preload("Files").First(invest, "id = ?", invest.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *InvestmentRepository) SoftDeleteInvestment(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, uid).Delete(&domain.Investment{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *InvestmentRepository) GetExpiredInvestment(ctx context.Context, olderThan time.Time) ([]domain.Investment, error) {
	var invests []domain.Investment
	if err := r.db.WithContext(ctx).Unscoped().Preload("Files").Where("deleted_at IS NOT NULL AND deleted_at < ?", olderThan).Find(&invests).Error; err != nil {
		return nil, err
	}

	return invests, nil
}

func (r *InvestmentRepository) HardDeleteInvestment(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("entity_id = ? AND entity_type = ?", id, "investment").Delete(&domain.FileAssociate{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Where("id = ?", id).Delete(&domain.Investment{}).Error; err != nil {
			return err
		}

		return nil
	})
}
