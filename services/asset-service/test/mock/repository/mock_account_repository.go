package mock

import (
	"context"
	"time"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) CreateAccount(ctx context.Context, acc *domain.Account) error {
	args := m.Called(ctx, acc)
	return args.Error(0)
}

func (m *MockAccountRepository) GetAccount(ctx context.Context, uid uuid.UUID) ([]*domain.Account, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) GetAccountByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Account, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) GetBatchAccountByIDs(ctx context.Context, ids []uuid.UUID) ([]*domain.Account, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) GetAccountByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) UpdateAccount(ctx context.Context, acc *domain.Account) (*domain.Account, error) {
	args := m.Called(ctx, acc)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) SoftDeleteAccount(ctx context.Context, id, uid uuid.UUID) error {
	args := m.Called(ctx, id, uid)
	return args.Error(0)
}

func (m *MockAccountRepository) GetExpiredAccounts(ctx context.Context, cutoff time.Time) ([]domain.Account, error) {
	args := m.Called(ctx, cutoff)
	return args.Get(0).([]domain.Account), args.Error(1)
}

func (m *MockAccountRepository) HardDeleteAccount(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
