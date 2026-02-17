package repository

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type BuildingRepository interface {
	CreateBuilding(ctx context.Context, item *domain.Building) error
	GetBuilding(ctx context.Context, uid uuid.UUID) ([]*domain.Building, error)
	GetBuildingByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Building, error)
	GetBatchBuildingByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Building, error)
	GetBuildingByID(ctx context.Context, id uuid.UUID) (*domain.Building, error)
	UpdateBuilding(ctx context.Context, item *domain.Building, addLandIDs, removeLandIDs, addInsIDs, removeInsIDs []uuid.UUID) (*domain.Building, error)
	SoftDeleteBuilding(ctx context.Context, id, uid uuid.UUID) error
	GetExpiredBuilding(ctx context.Context, olderThan time.Time) ([]domain.Building, error)
	HardDeleteBuilding(ctx context.Context, id uuid.UUID) error
}
