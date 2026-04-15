package mock

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockAssetRepository struct {
	mock.Mock
}

func (m *MockAssetRepository) GetAllAssetIDs(ctx context.Context, userID uuid.UUID) (map[string][]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(map[string][]string), args.Error(1)
}

func (m *MockAssetRepository) CheckExists(ctx context.Context, entityType string, id uuid.UUID, uid uuid.UUID) (string, bool, error) {
	args := m.Called(ctx, entityType, id, uid)
	return args.Get(0).(string), args.Get(1).(bool), args.Error(2)
}

func (m *MockAssetRepository) GetAllAssets(ctx context.Context, uid uuid.UUID) ([]domain.AssetSummary, []domain.AssetSummary, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]domain.AssetSummary), args.Get(1).([]domain.AssetSummary), args.Error(2)
}

func (m *MockAssetRepository) GetAssetCount(ctx context.Context, uid uuid.UUID) (int64, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAssetRepository) GetNetWorthOverview(ctx context.Context, uid uuid.UUID) (*domain.NetWorthOverview, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(*domain.NetWorthOverview), args.Error(1)
}
