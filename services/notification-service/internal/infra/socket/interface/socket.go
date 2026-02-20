package socket

type ISocketHub interface {
	Emit(userID string, data interface{})
	BroadcastToGroup(groupID string, data interface{})
}
