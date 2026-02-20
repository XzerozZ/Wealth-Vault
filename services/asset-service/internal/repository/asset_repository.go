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

func (r *AssetRepository) GetAllAssetIDs(ctx context.Context, userID uuid.UUID) (map[string][]string, error) {
	var results []domain.AssetIDResult

	query := `
		SELECT 'account' as type, id::text FROM accounts WHERE user_id = ? AND deleted_at IS NULL
		UNION ALL
		SELECT 'building' as type, id::text FROM buildings WHERE user_id = ? AND deleted_at IS NULL
		UNION ALL
		SELECT 'cash' as type, id::text FROM cash WHERE user_id = ? AND deleted_at IS NULL
		UNION ALL
		SELECT 'insurance' as type, id::text FROM insurances WHERE user_id = ? AND deleted_at IS NULL
		UNION ALL
		SELECT 'investment' as type, id::text FROM investments WHERE user_id = ? AND deleted_at IS NULL
		UNION ALL
		SELECT 'land' as type, id::text FROM lands WHERE user_id = ? AND deleted_at IS NULL
		UNION ALL
		SELECT 'liability' as type, id::text FROM liabilities WHERE user_id = ? AND deleted_at IS NULL
	`

	err := r.db.WithContext(ctx).Raw(query,
		userID, userID, userID, userID, userID, userID, userID,
	).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	assetMap := make(map[string][]string)
	for _, res := range results {
		assetMap[res.Type] = append(assetMap[res.Type], res.ID)
	}

	return assetMap, nil
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

func (r *AssetRepository) GetAllAssets(ctx context.Context, uid uuid.UUID) ([]domain.AssetSummary, error) {
	var assets []domain.AssetSummary

	query := `
		-- 1. Account
		SELECT id, 'account' as type, name, amount as value, created_at 
		FROM accounts WHERE user_id = ? AND deleted_at IS NULL

		UNION ALL

		-- 2. Building
		SELECT id, 'building' as type, name, amount as value, created_at 
		FROM buildings WHERE user_id = ? AND deleted_at IS NULL

		UNION ALL

		-- 3. Cash
		SELECT id, 'cash' as type, name, amount as value, created_at 
		FROM cashes WHERE user_id = ? AND deleted_at IS NULL

		UNION ALL

		-- 4. Insurance 
		SELECT id, 'insurance' as type, name, 0 as value, created_at 
		FROM insurances WHERE user_id = ? AND deleted_at IS NULL

		UNION ALL

		-- 5. Investment
		SELECT id, 'investment' as type, name, amount as value, created_at 
		FROM investments WHERE user_id = ? AND deleted_at IS NULL

		UNION ALL

		-- 6. Land
		SELECT id, 'land' as type, name, amount as value, created_at 
		FROM lands WHERE user_id = ? AND deleted_at IS NULL

		UNION ALL

		-- 7. Liability (Value = Principal)
		SELECT id, 'liability' as type, name, principal as value, created_at 
		FROM liabilities WHERE user_id = ? AND deleted_at IS NULL

		ORDER BY created_at DESC
	`

	err := r.db.WithContext(ctx).Raw(query, uid, uid, uid, uid, uid, uid, uid).Scan(&assets).Error

	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (r *AssetRepository) GetAssetCount(ctx context.Context, uid uuid.UUID) (int64, error) {
	var count int64

	query := `
		SELECT (
			-- 1. Account
			(SELECT COUNT(*) FROM accounts WHERE user_id = ? AND deleted_at IS NULL) +

			-- 2. Building
			(SELECT COUNT(*) FROM buildings WHERE user_id = ? AND deleted_at IS NULL) +

			-- 3. Cash
			(SELECT COUNT(*) FROM cashes WHERE user_id = ? AND deleted_at IS NULL) +

			-- 4. Investment
			(SELECT COUNT(*) FROM investments WHERE user_id = ? AND deleted_at IS NULL) +

			-- 5. Land
			(SELECT COUNT(*) FROM lands WHERE user_id = ? AND deleted_at IS NULL) +

			-- 6. Liability (นับเฉพาะที่ไม่ใช่ Expense)
			(SELECT COUNT(*) FROM liabilities WHERE user_id = ? AND deleted_at IS NULL AND type != 'Expense')
			
		) as total_count
	`

	err := r.db.WithContext(ctx).Raw(query, uid, uid, uid, uid, uid, uid).Scan(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *AssetRepository) GetNetWorthOverview(ctx context.Context, uid uuid.UUID) (*domain.NetWorthOverview, error) {
	var result domain.NetWorthOverview

	query := `
		SELECT 
			(
				COALESCE((SELECT SUM(amount) FROM accounts WHERE user_id = ? AND deleted_at IS NULL), 0) +
				COALESCE((SELECT SUM(amount) FROM buildings WHERE user_id = ? AND deleted_at IS NULL), 0) +
				COALESCE((SELECT SUM(amount) FROM cashes WHERE user_id = ? AND deleted_at IS NULL), 0) +
				COALESCE((SELECT SUM(amount) FROM investments WHERE user_id = ? AND deleted_at IS NULL), 0) +
				COALESCE((SELECT SUM(amount) FROM lands WHERE user_id = ? AND deleted_at IS NULL), 0)
			) as total_assets,

			COALESCE((SELECT SUM(principal) FROM liabilities WHERE user_id = ? AND deleted_at IS NULL AND type != 'Expense'), 0) as total_liabilities
	`

	err := r.db.WithContext(ctx).Raw(query, uid, uid, uid, uid, uid, uid).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}
