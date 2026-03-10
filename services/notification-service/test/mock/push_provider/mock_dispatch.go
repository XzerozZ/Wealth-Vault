package mock

import (
	"context"
	"wealth-vault/notification-service/internal/domain"

	push_provider "wealth-vault/notification-service/internal/infra/push_provider/interface"

	"github.com/stretchr/testify/mock"
)

type MockDispatcher struct {
	mock.Mock
}

func (m *MockDispatcher) SendToUser(ctx context.Context, tokens []domain.DeviceToken, payload push_provider.PushPayload) {
	m.Called(ctx, tokens, payload)
}
