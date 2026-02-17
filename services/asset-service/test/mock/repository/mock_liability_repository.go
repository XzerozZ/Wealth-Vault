package mock

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockLiabilityRepository struct {
	mock.Mock
}

func (m *MockLiabilityRepository) CreateLiability(ctx context.Context, asset *domain.Liability) error {
	args := m.Called(ctx, asset)
	return args.Error(0)
}

func (m *MockLiabilityRepository) GetLiability(ctx context.Context, uid uuid.UUID) ([]*domain.Liability, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]*domain.Liability), args.Error(1)
}

func (m *MockLiabilityRepository) GetLiabilityByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Liability, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*domain.Liability), args.Error(1)
}

func (m *MockLiabilityRepository) GetBatchLiabilityByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Liability, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*domain.Liability), args.Error(1)
}

func (m *MockLiabilityRepository) GetLiabilityByID(ctx context.Context, id uuid.UUID) (*domain.Liability, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Liability), args.Error(1)
}

func (m *MockLiabilityRepository) GetLiabilityByUserID(ctx context.Context, uid uuid.UUID) ([]*domain.Liability, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]*domain.Liability), args.Error(1)
}

func (m *MockLiabilityRepository) UpdateLiability(ctx context.Context, lia *domain.Liability) (*domain.Liability, error) {
	args := m.Called(ctx, lia)
	return args.Get(0).(*domain.Liability), args.Error(1)
}

func (m *MockLiabilityRepository) SoftDeleteLiability(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	args := m.Called(ctx, id, uid)
	return args.Error(0)
}

func (m *MockLiabilityRepository) GetExpiredLiability(ctx context.Context, olderThan time.Time) ([]domain.Liability, error) {
	args := m.Called(ctx, olderThan)
	return args.Get(0).([]domain.Liability), args.Error(1)
}

func (m *MockLiabilityRepository) HardDeleteLiability(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
