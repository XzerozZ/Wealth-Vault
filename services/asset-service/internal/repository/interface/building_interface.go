package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type BuildingRepository interface {
	CreateBuilding(ctx context.Context, item *domain.Building) error
	GetBuilding(ctx context.Context, uid uuid.UUID) ([]*domain.Building, error)
	GetBuildingByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Building, error)
	GetBuildingByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*domain.Building, error)
	UpdateBuilding(ctx context.Context, item *domain.Building, addLandIDs, removeLandIDs, addInsIDs, removeInsIDs []uuid.UUID) (*domain.Building, error)
	DeleteBuilding(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
}
