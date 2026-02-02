package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type InsuranceRepository interface {
	CreateInsurance(ctx context.Context, invest *domain.Insurance) error
	GetInsurance(ctx context.Context, uid uuid.UUID) ([]*domain.Insurance, error)
	GetInsuranceByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Insurance, error)
	GetInsuranceByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*domain.Insurance, error)
	UpdateInsurance(ctx context.Context, invest *domain.Insurance) (*domain.Insurance, error)
	DeleteInsurance(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
}
