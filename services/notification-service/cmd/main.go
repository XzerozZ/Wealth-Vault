package main

import (
	"log"
	"wealth-vault/notification-service/configs"
	"wealth-vault/notification-service/internal/delivery/http"
	"wealth-vault/notification-service/internal/delivery/worker"
	"wealth-vault/notification-service/internal/infra/nats"
	"wealth-vault/notification-service/internal/infra/socket"
	"wealth-vault/notification-service/internal/repository"
	"wealth-vault/notification-service/internal/usecase"
	"wealth-vault/notification-service/pkg/database"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := configs.LoadConfigs()
	database.InitDB(cfg.PostgreSQL)
	db := database.GetDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
	}

	nc, _ := nats.NewNATSConnection(cfg.NATS.Host, cfg.NATS.Port)
	hub := socket.NewSocketHub()

	repo := repository.NewNotificationRepository(db)
	uc := usecase.NewNotificationUsecase(repo, hub)

	worker.StartConsumer(nc, uc)

	app := fiber.New()
	handler := http.NewHandler(hub, uc)

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws", websocket.New(handler.WebSocketEndpoint))

	app.Get("/notifications", handler.GetNotifications)

	app.Listen(":" + cfg.APP.Port)
}
