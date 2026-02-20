package mock

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockBuildingRepository struct {
	mock.Mock
}

func (m *MockBuildingRepository) CreateBuilding(ctx context.Context, item *domain.Building) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockBuildingRepository) GetBuilding(ctx context.Context, uid uuid.UUID) ([]*domain.Building, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]*domain.Building), args.Error(1)
}

func (m *MockBuildingRepository) GetBuildingByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Building, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*domain.Building), args.Error(1)
}

func (m *MockBuildingRepository) GetBatchBuildingByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Building, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*domain.Building), args.Error(1)
}

func (m *MockBuildingRepository) GetBuildingByID(ctx context.Context, id uuid.UUID) (*domain.Building, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Building), args.Error(1)
}

func (m *MockBuildingRepository) UpdateBuilding(ctx context.Context, item *domain.Building, addLandIDs, removeLandIDs, addInsIDs, removeInsIDs []uuid.UUID) (*domain.Building, error) {
	args := m.Called(ctx, item, addLandIDs, removeLandIDs, addInsIDs, removeInsIDs)
	return args.Get(0).(*domain.Building), args.Error(1)
}

func (m *MockBuildingRepository) SoftDeleteBuilding(ctx context.Context, id, uid uuid.UUID) error {
	args := m.Called(ctx, id, uid)
	return args.Error(0)
}

func (m *MockBuildingRepository) GetExpiredBuilding(ctx context.Context, olderThan time.Time) ([]domain.Building, error) {
	args := m.Called(ctx, olderThan)
	return args.Get(0).([]domain.Building), args.Error(1)
}

func (m *MockBuildingRepository) HardDeleteBuilding(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
