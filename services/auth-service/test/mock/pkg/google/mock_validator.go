package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"google.golang.org/api/idtoken"
)

type MockGoogleValidator struct {
	mock.Mock
}

func (m *MockGoogleValidator) Validate(ctx context.Context, token string) (*idtoken.Payload, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*idtoken.Payload), args.Error(1)
}
