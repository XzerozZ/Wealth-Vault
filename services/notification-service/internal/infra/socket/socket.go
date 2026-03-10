package socket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type client struct {
	connID string
	conn   *websocket.Conn
	mu     sync.Mutex
}

type SocketHub struct {
	clients map[string]map[string]*client
	rooms   map[string]map[string]bool
	mu      sync.RWMutex
}

func NewSocketHub() *SocketHub {
	return &SocketHub{
		clients: make(map[string]map[string]*client),
		rooms:   make(map[string]map[string]bool),
	}
}

func (h *SocketHub) Register(userID string, conn *websocket.Conn) string {
	connID := uuid.New().String()

	h.mu.Lock()
	if h.clients[userID] == nil {
		h.clients[userID] = make(map[string]*client)
	}
	h.clients[userID][connID] = &client{connID: connID, conn: conn}
	total := len(h.clients[userID])
	h.mu.Unlock()

	log.Printf("🔌 User %s connected | connID: %.8s | devices online: %d", userID, connID, total)
	return connID
}

func (h *SocketHub) Unregister(userID, connID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conns, ok := h.clients[userID]
	if !ok {
		return
	}

	if c, exists := conns[connID]; exists {
		c.conn.Close()
		delete(conns, connID)
	}

	if len(conns) == 0 {
		delete(h.clients, userID)
		for groupID, members := range h.rooms {
			delete(members, userID)
			if len(members) == 0 {
				delete(h.rooms, groupID)
			}
		}
		log.Printf("❌ User %s fully offline (all devices disconnected)", userID)
	} else {
		log.Printf("📴 User %s device disconnected | connID: %.8s | remaining: %d", userID, connID, len(conns))
	}
}

func (h *SocketHub) JoinGroup(userID, groupID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[groupID] == nil {
		h.rooms[groupID] = make(map[string]bool)
	}

	h.rooms[groupID][userID] = true
	log.Printf("➕ User %s joined group %s", userID, groupID)
}

func (h *SocketHub) LeaveGroup(userID, groupID string) {
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
	snapshot := make(map[string]*client, len(h.clients[userID]))
	for id, c := range h.clients[userID] {
		snapshot[id] = c
	}

	h.mu.RUnlock()
	for connID, c := range snapshot {
		if err := h.send(c, data); err != nil {
			log.Printf("⚠️  Write error to %s [%.8s]: %v", userID, connID, err)
			h.Unregister(userID, connID)
		}
	}
}

func (h *SocketHub) BroadcastToGroup(groupID string, data interface{}) {
	h.mu.RLock()
	members, ok := h.rooms[groupID]
	if !ok {
		h.mu.RUnlock()
		return
	}

	userIDs := make([]string, 0, len(members))
	for uid := range members {
		userIDs = append(userIDs, uid)
	}

	h.mu.RUnlock()

	for _, uid := range userIDs {
		h.Emit(uid, data)
	}
}

func (h *SocketHub) BroadcastToGroupExcept(groupID, senderID string, data interface{}) {
	h.mu.RLock()
	members, ok := h.rooms[groupID]
	if !ok {
		h.mu.RUnlock()
		return
	}

	userIDs := make([]string, 0, len(members))
	for uid := range members {
		if uid != senderID {
			userIDs = append(userIDs, uid)
		}
	}

	h.mu.RUnlock()

	for _, uid := range userIDs {
		h.Emit(uid, data)
	}
}

func (h *SocketHub) IsOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[userID]) > 0
}

func (h *SocketHub) send(c *client, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(websocket.TextMessage, b)
}
