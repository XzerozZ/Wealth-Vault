package http

import (
	"wealth-vault/notification-service/internal/infra/socket"
	"wealth-vault/notification-service/internal/usecase"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	hub *socket.SocketHub
	uc  *usecase.NotificationUsecase
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
	defer func() { h.hub.Unregister(userID); c.Close() }()

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}

func (h *Handler) GetNotifications(c *fiber.Ctx) error {
	userID := c.Get("X-User-ID")
	if userID == "" {
		return c.SendStatus(401)
	}

	history, err := h.uc.GetHistory(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(history)
}
