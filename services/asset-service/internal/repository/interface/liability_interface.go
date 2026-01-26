package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type LiabilityRepository interface {
	CreateLiability(ctx context.Context, asset *domain.Liability) error
	GetLiability(ctx context.Context, uid uuid.UUID) ([]*domain.Liability, error)
	GetLiabilityByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*domain.Liability, error)
	UpdateLiability(ctx context.Context, lia *domain.Liability) (*domain.Liability, error)
	DeleteLiability(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
}
