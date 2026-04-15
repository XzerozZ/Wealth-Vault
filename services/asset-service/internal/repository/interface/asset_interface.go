package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type AssetRepository interface {
	GetAllAssetIDs(ctx context.Context, userID uuid.UUID) (map[string][]string, error)
	CheckExists(ctx context.Context, entityType string, id uuid.UUID, uid uuid.UUID) (string, bool, error)
	GetAllAssets(ctx context.Context, uid uuid.UUID) ([]domain.AssetSummary, []domain.AssetSummary, error)
	GetAssetCount(ctx context.Context, uid uuid.UUID) (int64, error)
	GetNetWorthOverview(ctx context.Context, uid uuid.UUID) (*domain.NetWorthOverview, error)
}
