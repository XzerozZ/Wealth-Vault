package mock

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/stretchr/testify/mock"
)

type MockSocketHub struct {
	mock.Mock
}

func (m *MockSocketHub) Register(userID string, conn *websocket.Conn) string {
	args := m.Called(userID)
	return args.String(0)
}

func (m *MockSocketHub) Unregister(userID string, connID string) {
	m.Called(userID, connID)
}

func (m *MockSocketHub) Emit(userID string, data interface{}) {
	m.Called(userID, data)
}

func (m *MockSocketHub) BroadcastToGroup(groupID string, data interface{}) {
	m.Called(groupID, data)
}

func (m *MockSocketHub) BroadcastToGroupExcept(groupID string, senderID string, data interface{}) {
	m.Called(groupID, senderID, data)
}

func (m *MockSocketHub) JoinGroup(userID string, groupID string) {
	m.Called(userID, groupID)
}

func (m *MockSocketHub) LeaveGroup(userID string, groupID string) {
	m.Called(userID, groupID)
}

func (m *MockSocketHub) IsOnline(userID string) bool {
	args := m.Called(userID)
	return args.Bool(0)
}
