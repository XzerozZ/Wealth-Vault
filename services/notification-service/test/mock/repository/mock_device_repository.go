package mock

import (
	"context"
	"wealth-vault/notification-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockDeviceRepository struct {
	mock.Mock
}

func (m *MockDeviceRepository) RegisterDevice(ctx context.Context, req *domain.DeviceToken) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockDeviceRepository) GetActiveTokens(ctx context.Context, userID uuid.UUID) ([]domain.DeviceToken, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) != nil {
		return args.Get(0).([]domain.DeviceToken), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDeviceRepository) UnregisterDevice(ctx context.Context, userID uuid.UUID, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

func (m *MockDeviceRepository) MarkTokenInactive(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
