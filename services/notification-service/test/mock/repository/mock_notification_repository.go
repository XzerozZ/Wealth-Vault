package mock

import (
	"context"
	"wealth-vault/notification-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) CreateNotification(ctx context.Context, item *domain.Notification) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockNotificationRepository) GetByReceiver(ctx context.Context, receiverID uuid.UUID) ([]domain.Notification, error) {
	args := m.Called(ctx, receiverID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]domain.Notification), args.Error(1)
}
