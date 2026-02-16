package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type AssetRepository interface {
	CheckExists(ctx context.Context, entityType string, id uuid.UUID, uid uuid.UUID) (bool, error)
	GetAllAssets(ctx context.Context, uid uuid.UUID) ([]domain.AssetSummary, error)
	GetAssetCount(ctx context.Context, uid uuid.UUID) (int64, error)
	GetNetWorthOverview(ctx context.Context, uid uuid.UUID) (*domain.NetWorthOverview, error)
}
