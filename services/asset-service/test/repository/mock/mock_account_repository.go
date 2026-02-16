package mock

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepo struct {
	mock.Mock
}

func (m *MockAccountRepo) CreateAccount(ctx context.Context, acc *domain.Account) error {
	args := m.Called(ctx, acc)
	return args.Error(0)
}

func (m *MockAccountRepo) GetAccount(ctx context.Context, uid uuid.UUID) ([]*domain.Account, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).([]*domain.Account), args.Error(1)
}

func (m *MockAccountRepo) GetAccountByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepo) UpdateAccount(ctx context.Context, acc *domain.Account) (*domain.Account, error) {
	args := m.Called(ctx, acc)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepo) SoftDeleteAccount(ctx context.Context, id, uid uuid.UUID) error {
	args := m.Called(ctx, id, uid)
	return args.Error(0)
}
