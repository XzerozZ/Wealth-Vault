package main

import (
	"log"
	grpcclient "wealth-vault/notification-service/client"
	"wealth-vault/notification-service/configs"
	"wealth-vault/notification-service/internal/delivery/http"
	"wealth-vault/notification-service/internal/delivery/worker"
	"wealth-vault/notification-service/internal/infra/database"
	line "wealth-vault/notification-service/internal/infra/line"
	"wealth-vault/notification-service/internal/infra/nats"
	push_provider "wealth-vault/notification-service/internal/infra/push_provider"
	socket "wealth-vault/notification-service/internal/infra/socket"
	"wealth-vault/notification-service/internal/repository"
	"wealth-vault/notification-service/internal/usecase"

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

	fcmClient, err := push_provider.NewFCMProvider(cfg.FCM.CredentialsFile)
	if err != nil {
		log.Fatalf("Failed to initialize FCM: %v", err)
	}

	authClient, err := grpcclient.NewAuthClient(cfg.AuthService.Host, cfg.AuthService.Port)
	if err != nil {
		log.Fatal("auth service:", err)
	}

	lineAPIClient := line.NewLineClient(cfg.Line.Token)

	// ------ repository ------
	notiRepo := repository.NewNotificationRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	dispatcher := push_provider.NewDispatcher(fcmClient, nil, deviceRepo)

	deviceUC := usecase.NewDeviceUsecase(deviceRepo)
	notiUC := usecase.NewNotificationUsecase(notiRepo, deviceRepo, hub, dispatcher, lineAPIClient, authClient)

	worker.StartConsumer(nc, notiUC)

	app := fiber.New()

	notiHandler := http.NewHandler(hub, notiUC)
	deviceHandler := http.NewDeviceHandler(deviceUC)

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws", websocket.New(notiHandler.WebSocketEndpoint))

	notiGroup := app.Group("/notifications")
	notiGroup.Get("/", notiHandler.GetNotifications)
	notiGroup.Put("/read-all", notiHandler.MarkAllAsRead)
	notiGroup.Put("/:id/read", notiHandler.MarkAsRead)

	deviceGroup := app.Group("/devices")
	deviceGroup.Post("/register", deviceHandler.RegisterDevice)
	deviceGroup.Post("/unregister", deviceHandler.UnregisterDevice)
	deviceGroup.Get("/", deviceHandler.GetDevices)

	app.Listen(":" + cfg.APP.Port)
}
