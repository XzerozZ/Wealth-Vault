package mock

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockInvestmentRepository struct {
	mock.Mock
}

func (m *MockInvestmentRepository) CreateInvestment(ctx context.Context, invest *domain.Investment) error {
	args := m.Called(ctx, invest)
	return args.Error(0)
}

func (m *MockInvestmentRepository) GetInvestment(ctx context.Context, uid uuid.UUID) ([]*domain.Investment, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]*domain.Investment), args.Error(1)
}

func (m *MockInvestmentRepository) GetInvestmentByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Investment, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*domain.Investment), args.Error(1)
}

func (m *MockInvestmentRepository) GetBatchInvestmentByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Investment, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*domain.Investment), args.Error(1)
}

func (m *MockInvestmentRepository) GetInvestmentByID(ctx context.Context, id uuid.UUID) (*domain.Investment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Investment), args.Error(1)
}

func (m *MockInvestmentRepository) GetInvestmentByUserID(ctx context.Context, uid uuid.UUID) ([]*domain.Investment, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]*domain.Investment), args.Error(1)
}

func (m *MockInvestmentRepository) UpdateInvestment(ctx context.Context, invest *domain.Investment) (*domain.Investment, error) {
	args := m.Called(ctx, invest)
	return args.Get(0).(*domain.Investment), args.Error(1)
}

func (m *MockInvestmentRepository) SoftDeleteInvestment(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	args := m.Called(ctx, id, uid)
	return args.Error(0)
}

func (m *MockInvestmentRepository) GetExpiredInvestment(ctx context.Context, olderThan time.Time) ([]domain.Investment, error) {
	args := m.Called(ctx, olderThan)
	return args.Get(0).([]domain.Investment), args.Error(1)
}

func (m *MockInvestmentRepository) HardDeleteInvestment(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
