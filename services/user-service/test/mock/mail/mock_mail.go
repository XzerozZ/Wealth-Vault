package mock

import (
	"context"
	"wealth-vault/user-service/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockMailClient struct {
	mock.Mock
}

func (m *MockMailClient) SendShareInvitation(ctx context.Context, req domain.SendEmailRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
