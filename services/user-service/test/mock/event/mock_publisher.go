package mock

import (
	"github.com/stretchr/testify/mock"
)

type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) Publish(subject string, v any) error {
	args := m.Called(subject, v)
	return args.Error(0)
}
