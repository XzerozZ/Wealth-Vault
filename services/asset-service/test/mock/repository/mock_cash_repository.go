package mock

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockCashRepository struct {
	mock.Mock
}

func (m *MockCashRepository) CreateCash(ctx context.Context, cash *domain.Cash) error {
	args := m.Called(ctx, cash)
	return args.Error(0)
}

func (m *MockCashRepository) GetCash(ctx context.Context, uid uuid.UUID) ([]*domain.Cash, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]*domain.Cash), args.Error(1)
}

func (m *MockCashRepository) GetCashByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Cash, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*domain.Cash), args.Error(1)
}

func (m *MockCashRepository) GetBatchCashByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Cash, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*domain.Cash), args.Error(1)
}

func (m *MockCashRepository) GetCashByID(ctx context.Context, id uuid.UUID) (*domain.Cash, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Cash), args.Error(1)
}

func (m *MockCashRepository) UpdateCash(ctx context.Context, cash *domain.Cash) (*domain.Cash, error) {
	args := m.Called(ctx, cash)
	return args.Get(0).(*domain.Cash), args.Error(1)
}

func (m *MockCashRepository) SoftDeleteCash(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	args := m.Called(ctx, id, uid)
	return args.Error(0)
}

func (m *MockCashRepository) HardDeleteCash(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCashRepository) GetExpiredCash(ctx context.Context, olderThan time.Time) ([]domain.Cash, error) {
	args := m.Called(ctx, olderThan)
	return args.Get(0).([]domain.Cash), args.Error(1)
}
