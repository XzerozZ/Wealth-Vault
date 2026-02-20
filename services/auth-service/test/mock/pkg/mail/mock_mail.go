package mock

import (
	"context"
	"wealth-vault/auth-service/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockNotificationClient struct {
	mock.Mock
}

func (m *MockNotificationClient) SendOTP(ctx context.Context, req domain.SendEmailRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
