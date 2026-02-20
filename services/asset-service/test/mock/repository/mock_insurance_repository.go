package mock

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockInsuranceRepository struct {
	mock.Mock
}

func (m *MockInsuranceRepository) CreateInsurance(ctx context.Context, invest *domain.Insurance) error {
	args := m.Called(ctx, invest)
	return args.Error(0)
}

func (m *MockInsuranceRepository) GetInsurance(ctx context.Context, uid uuid.UUID) ([]*domain.Insurance, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]*domain.Insurance), args.Error(1)
}

func (m *MockInsuranceRepository) GetInsuranceByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Insurance, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*domain.Insurance), args.Error(1)
}

func (m *MockInsuranceRepository) GetBatchInsuranceByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Insurance, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*domain.Insurance), args.Error(1)
}

func (m *MockInsuranceRepository) GetInsuranceByID(ctx context.Context, id uuid.UUID) (*domain.Insurance, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Insurance), args.Error(1)
}

func (m *MockInsuranceRepository) UpdateInsurance(ctx context.Context, invest *domain.Insurance) (*domain.Insurance, error) {
	args := m.Called(ctx, invest)
	return args.Get(0).(*domain.Insurance), args.Error(1)
}

func (m *MockInsuranceRepository) SoftDeleteInsurances(ctx context.Context, id uuid.UUID, uid uuid.UUID) error {
	args := m.Called(ctx, id, uid)
	return args.Error(0)
}

func (m *MockInsuranceRepository) GetExpiredInsurances(ctx context.Context, olderThan time.Time) ([]domain.Insurance, error) {
	args := m.Called(ctx, olderThan)
	return args.Get(0).([]domain.Insurance), args.Error(1)
}

func (m *MockInsuranceRepository) HardDeleteInsurances(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockInsuranceRepository) GetExpiringInsurances(ctx context.Context, days int) ([]*domain.Insurance, error) {
	args := m.Called(ctx, days)
	return args.Get(0).([]*domain.Insurance), args.Error(1)
}
