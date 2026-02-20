package mock

import "github.com/stretchr/testify/mock"

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(topic string, evt interface{}) error {
	args := m.Called(topic, evt)
	return args.Error(0)
}
