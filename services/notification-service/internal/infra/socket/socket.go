package socket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type SocketHub struct {
	clients map[string]*websocket.Conn
	rooms   map[string]map[string]bool
	mu      sync.RWMutex
}

func NewSocketHub() *SocketHub {
	return &SocketHub{
		clients: make(map[string]*websocket.Conn),
		rooms:   make(map[string]map[string]bool),
	}
}

func (h *SocketHub) Register(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = conn
	log.Printf("🔌 User %s Connected (Online)", userID)
}

func (h *SocketHub) Unregister(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[userID]; ok {
		delete(h.clients, userID)
	}

	for groupID, members := range h.rooms {
		if _, ok := members[userID]; ok {
			delete(members, userID)
			if len(members) == 0 {
				delete(h.rooms, groupID)
			}
		}
	}
	log.Printf("❌ User %s Disconnected (Cleaned up)", userID)
}

func (h *SocketHub) JoinGroup(userID string, groupID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[groupID] == nil {
		h.rooms[groupID] = make(map[string]bool)
	}
	h.rooms[groupID][userID] = true
	log.Printf("➕ User %s JOINED Group %s", userID, groupID)
}

func (h *SocketHub) LeaveGroup(userID string, groupID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if members, ok := h.rooms[groupID]; ok {
		delete(members, userID)
		log.Printf("➖ User %s LEFT Group %s", userID, groupID)
		if len(members) == 0 {
			delete(h.rooms, groupID)
		}
	}
}

func (h *SocketHub) Emit(userID string, data interface{}) {
	h.mu.RLock()
	conn, ok := h.clients[userID]
	h.mu.RUnlock()

	if ok {
		h.send(userID, conn, data)
	}
}

func (h *SocketHub) BroadcastToGroup(groupID string, data interface{}) {
	h.mu.RLock()
	members, ok := h.rooms[groupID]
	h.mu.RUnlock()

	if !ok {
		return
	}

	for memberID := range members {
		h.mu.RLock()
		conn, ok := h.clients[memberID]
		h.mu.RUnlock()
		if ok {
			h.send(memberID, conn, data)
		}
	}
}

func (h *SocketHub) send(userID string, conn *websocket.Conn, data interface{}) {
	msgBytes, _ := json.Marshal(data)
	if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		conn.Close()
		go h.Unregister(userID)
	}
}
