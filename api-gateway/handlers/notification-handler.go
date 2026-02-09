package handlers

import (
	"fmt"
	"log"
	"strings"

	"wealth-vault/api-gateway/configs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"

	fasthttpws "github.com/fasthttp/websocket"
	"github.com/gofiber/websocket/v2"
)

type NotificationHandler struct {
	cfg *configs.Configs
}

func NewNotificationHandler(cfg *configs.Configs) *NotificationHandler {
	return &NotificationHandler{
		cfg: cfg,
	}
}

func (h *NotificationHandler) ProxyWebSocket(c *websocket.Conn) {
	userID, _ := c.Locals("user_id").(string)
	if userID == "" {
		return
	}

	targetURL := fmt.Sprintf("ws://%s:%s/ws?user_id=%s",
		h.cfg.NotiService.Host, h.cfg.NotiService.Port, userID)

	h.bridgeWS(c, targetURL, userID)
}

func (h *NotificationHandler) bridgeWS(c *websocket.Conn, targetURL string, userID string) {
	dialer := fasthttpws.DefaultDialer
	upstream, _, err := dialer.Dial(targetURL, nil)
	if err != nil {
		log.Printf("❌ Gateway Dial Error [%s]: %v", userID, err)
		return
	}
	defer upstream.Close()

	log.Printf("🔗 Connected: %s", userID)

	errChan := make(chan error, 1)

	go func() {
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			if err := upstream.WriteMessage(mt, msg); err != nil {
				errChan <- err
				return
			}
		}
	}()

	go func() {
		for {
			mt, msg, err := upstream.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			if err := c.WriteMessage(mt, msg); err != nil {
				errChan <- err
				return
			}
		}
	}()

	<-errChan
	log.Printf("🔌 Disconnected: %s", userID)
}

func (h *NotificationHandler) ProxyAPI(c *fiber.Ctx) error {
	baseURL := fmt.Sprintf("http://%s:%s", h.cfg.NotiService.Host, h.cfg.NotiService.Port)

	path := strings.TrimPrefix(c.OriginalURL(), "/api")
	targetURL := baseURL + path

	userID := c.Locals("user_id").(string)
	c.Request().Header.Set("X-User-ID", userID)

	return proxy.Do(c, targetURL)
}
