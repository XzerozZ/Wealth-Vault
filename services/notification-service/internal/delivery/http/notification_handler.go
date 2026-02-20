package http

import (
	"encoding/json"
	"log"
	"wealth-vault/notification-service/internal/infra/socket"
	"wealth-vault/notification-service/internal/usecase"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	hub *socket.SocketHub
	uc  *usecase.NotificationUsecase
}

type ClientCommand struct {
	Action  string `json:"action"`
	GroupID string `json:"group_id"`
}

func NewHandler(hub *socket.SocketHub, uc *usecase.NotificationUsecase) *Handler {
	return &Handler{hub: hub, uc: uc}
}

func (h *Handler) WebSocketEndpoint(c *websocket.Conn) {
	userID := c.Query("user_id")
	if userID == "" {
		c.Close()
		return
	}

	h.hub.Register(userID, c)
	defer func() {
		h.hub.Unregister(userID)
		c.Close()
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}

		var cmd ClientCommand
		if err := json.Unmarshal(msg, &cmd); err != nil {
			log.Printf("⚠️ Invalid Command from %s: %v", userID, err)
			continue
		}

		switch cmd.Action {
		case "JOIN":
			if cmd.GroupID != "" {
				h.hub.JoinGroup(userID, cmd.GroupID)
			}
		case "LEAVE":
			if cmd.GroupID != "" {
				h.hub.LeaveGroup(userID, cmd.GroupID)
			}
		default:
			log.Printf("⚠️ Unknown Action: %s", cmd.Action)
		}
	}
}

func (h *Handler) GetNotifications(c *fiber.Ctx) error {
	userID := c.Get("X-User-ID")
	if userID == "" {
		return c.SendStatus(401)
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid user id"})
	}

	history, err := h.uc.GetHistory(c.Context(), uid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(history)
}
