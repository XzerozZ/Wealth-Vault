package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type AssetRepository interface {
	CreateAsset(ctx context.Context, asset *domain.Asset) error
	GetAsset(ctx context.Context, uid uuid.UUID) ([]*domain.Asset, error)
	GetAssetByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*domain.Asset, error)
	UpdateAsset(ctx context.Context, asset *domain.Asset, mask []string) (*domain.Asset, error)
	DeleteAsset(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
}
