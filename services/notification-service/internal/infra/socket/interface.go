package socket

import "github.com/gofiber/contrib/websocket"

type ISocketHub interface {
	Register(userID string, conn *websocket.Conn) string
	Unregister(userID string, connID string)
	Emit(userID string, data interface{})
	BroadcastToGroup(groupID string, data interface{})
	BroadcastToGroupExcept(groupID string, senderID string, data interface{})
	JoinGroup(userID string, groupID string)
	LeaveGroup(userID string, groupID string)
	IsOnline(userID string) bool
}
