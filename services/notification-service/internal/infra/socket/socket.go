package socket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type SocketHub struct {
	clients sync.Map
}

func NewSocketHub() *SocketHub { return &SocketHub{} }

func (h *SocketHub) Register(userID string, conn *websocket.Conn) {
	h.clients.Store(userID, conn)
}

func (h *SocketHub) Unregister(userID string) {
	h.clients.Delete(userID)
}

func (h *SocketHub) Emit(userID string, data interface{}) {
	val, ok := h.clients.Load(userID)
	if !ok {
		log.Printf("⚠️ User %s is OFFLINE (No Socket Connection)", userID)
		return
	}

	log.Printf("✅ User %s is ONLINE. Sending data...", userID)

	conn := val.(*websocket.Conn)
	msgBytes, _ := json.Marshal(data)
	if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		conn.Close()
		h.Unregister(userID)
	}
}
