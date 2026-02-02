package handlers

import (
	"fmt"
	"log"
	"strings"

	"wealth-vault/api-gateway/configs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"

	// 1. อันนี้ของ Fiber (เอาไว้คุยกับ Frontend)
	"github.com/gofiber/websocket/v2"

	// 2. ✅ เพิ่มอันนี้! (เอาไว้โทรหา Backend Notification Service)
	fasthttpws "github.com/fasthttp/websocket"
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
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		log.Println("❌ user_id not found")
		return
	}

	targetURL := fmt.Sprintf(
		"ws://%s:%s/ws?user_id=%s",
		h.cfg.NotiService.Host,
		h.cfg.NotiService.Port,
		userID,
	)

	upstream, _, err := fasthttpws.DefaultDialer.Dial(targetURL, nil)
	if err != nil {
		log.Println("❌ cannot connect to notification service:", err)
		return
	}
	defer upstream.Close()

	log.Println("🔗 WS bridge connected:", userID)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				return
			}
			if err := upstream.WriteMessage(mt, msg); err != nil {
				return
			}
		}
	}()

	for {
		mt, msg, err := upstream.ReadMessage()
		if err != nil {
			break
		}

		log.Printf("⬅️ Gateway received from Service: %s", string(msg))

		if err := c.WriteMessage(mt, msg); err != nil {
			break
		}
	}
}

func (h *NotificationHandler) ProxyAPI(c *fiber.Ctx) error {
	baseURL := fmt.Sprintf("http://%s:%s", h.cfg.NotiService.Host, h.cfg.NotiService.Port)

	path := strings.TrimPrefix(c.OriginalURL(), "/api")
	targetURL := baseURL + path

	userID := c.Locals("user_id").(string)
	c.Request().Header.Set("X-User-ID", userID)

	return proxy.Do(c, targetURL)
}
