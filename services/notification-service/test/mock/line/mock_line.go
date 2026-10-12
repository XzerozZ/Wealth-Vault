package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockLineClient struct {
	mock.Mock
}

func (m *MockLineClient) SendTextMessage(ctx context.Context, to string, text string, retryKey string) error {
	args := m.Called(ctx, to, text, retryKey)
	return args.Error(0)
}
