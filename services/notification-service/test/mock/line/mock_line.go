package mock

import "github.com/stretchr/testify/mock"

type MockLineClient struct {
	mock.Mock
}

func (m *MockLineClient) SendTextMessage(to string, text string) error {
	args := m.Called(to, text)
	return args.Error(0)
}
