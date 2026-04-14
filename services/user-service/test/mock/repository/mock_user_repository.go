package mock

import (
	"context"
	"wealth-vault/user-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)

	var user *domain.User
	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User)
	}

	return user, args.Error(1)
}

func (m *MockUserRepository) GetUsersByEmail(ctx context.Context, email string) ([]*domain.User, error) {
	args := m.Called(ctx, email)

	var user []*domain.User
	if args.Get(0) != nil {
		user = args.Get(0).([]*domain.User)
	}

	return user, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *domain.User, mask []string) (*domain.User, error) {
	args := m.Called(ctx, user, mask)

	var updated *domain.User
	if args.Get(0) != nil {
		updated = args.Get(0).(*domain.User)
	}

	return updated, args.Error(1)
}
