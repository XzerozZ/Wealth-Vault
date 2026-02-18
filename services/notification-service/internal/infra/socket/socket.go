package socket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type client struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

type SocketHub struct {
	clients map[string]*client
	rooms   map[string]map[string]bool
	mu      sync.RWMutex
}

func NewSocketHub() *SocketHub {
	return &SocketHub{
		clients: make(map[string]*client),
		rooms:   make(map[string]map[string]bool),
	}
}

func (h *SocketHub) Register(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[userID] = &client{conn: conn}
	h.mu.Unlock()

	log.Printf("🔌 User %s Connected (Online)", userID)
}

func (h *SocketHub) Unregister(userID string) {
	h.mu.Lock()

	c, ok := h.clients[userID]
	if ok {
		delete(h.clients, userID)
		c.conn.Close()
	}

	for groupID, members := range h.rooms {
		if _, ok := members[userID]; ok {
			delete(members, userID)
			if len(members) == 0 {
				delete(h.rooms, groupID)
			}
		}
	}

	h.mu.Unlock()

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
		if len(members) == 0 {
			delete(h.rooms, groupID)
		}
	}

	log.Printf("➖ User %s LEFT Group %s", userID, groupID)
}

func (h *SocketHub) Emit(userID string, data interface{}) {
	h.mu.RLock()
	c, ok := h.clients[userID]
	h.mu.RUnlock()

	if ok {
		h.send(userID, c, data)
	}
}

func (h *SocketHub) BroadcastToGroup(groupID string, data interface{}) {
	h.mu.RLock()
	membersMap, ok := h.rooms[groupID]
	if !ok {
		h.mu.RUnlock()
		return
	}

	memberIDs := make([]string, 0, len(membersMap))
	for id := range membersMap {
		memberIDs = append(memberIDs, id)
	}
	h.mu.RUnlock()

	for _, id := range memberIDs {
		h.mu.RLock()
		c, ok := h.clients[id]
		h.mu.RUnlock()

		if ok {
			h.send(id, c, data)
		}
	}
}

func (h *SocketHub) send(userID string, c *client, data interface{}) {
	msgBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("❌ Marshal Error: %v", err)
		return
	}

	c.mu.Lock()
	err = c.conn.WriteMessage(websocket.TextMessage, msgBytes)
	c.mu.Unlock()

	if err != nil {
		log.Printf("⚠️ Write Error to %s: %v", userID, err)
		h.Unregister(userID)
	}
}
