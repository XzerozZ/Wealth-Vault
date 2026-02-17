package mock

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockLandRepository struct {
	mock.Mock
}

func (m *MockLandRepository) CreateLand(ctx context.Context, item *domain.Land) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockLandRepository) GetLand(ctx context.Context, uid uuid.UUID) ([]*domain.Land, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]*domain.Land), args.Error(1)
}

func (m *MockLandRepository) GetLandByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Land, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*domain.Land), args.Error(1)
}

func (m *MockLandRepository) GetBatchLandByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Land, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*domain.Land), args.Error(1)
}

func (m *MockLandRepository) GetLandByID(ctx context.Context, id uuid.UUID) (*domain.Land, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Land), args.Error(1)
}

func (m *MockLandRepository) GetLandByUserID(ctx context.Context, uid uuid.UUID) ([]*domain.Land, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]*domain.Land), args.Error(1)
}

func (m *MockLandRepository) UpdateLand(ctx context.Context, item *domain.Land, addBuildIDs []uuid.UUID, removeBuildIDs []uuid.UUID) (*domain.Land, error) {
	args := m.Called(ctx, item, addBuildIDs, removeBuildIDs)
	return args.Get(0).(*domain.Land), args.Error(1)
}

func (m *MockLandRepository) SoftDeleteLand(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	args := m.Called(ctx, id, uid)
	return args.Error(0)
}

func (m *MockLandRepository) GetExpiredLand(ctx context.Context, olderThan time.Time) ([]domain.Land, error) {
	args := m.Called(ctx, olderThan)
	return args.Get(0).([]domain.Land), args.Error(1)
}

func (m *MockLandRepository) HardDeleteLand(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
