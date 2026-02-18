package mock

import "github.com/stretchr/testify/mock"

type MockSocketHub struct {
	mock.Mock
}

func (m *MockSocketHub) Emit(userID string, data interface{}) {
	m.Called(userID, data)
}

func (m *MockSocketHub) BroadcastToGroup(groupID string, data interface{}) {
	m.Called(groupID, data)
}
