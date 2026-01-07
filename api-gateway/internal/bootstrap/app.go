package bootstrap

import (
	"log"

	grpcclient "wealth-vault/api-gateway/client"
	"wealth-vault/api-gateway/configs"
	"wealth-vault/api-gateway/handlers"
	"wealth-vault/api-gateway/routes"

	"github.com/gofiber/fiber/v2"
)

func InitApp(cfg *configs.Configs) *fiber.App {
	app := fiber.New()

	// ---------- gRPC clients ----------
	userClient, err := grpcclient.NewUserClient(cfg.UserService.Host, cfg.UserService.Port)
	if err != nil {
		log.Fatal("user service:", err)
	}

	authClient, err := grpcclient.NewAuthClient(cfg.AuthService.Host, cfg.AuthService.Port)
	if err != nil {
		log.Fatal("auth service:", err)
	}

	// ---------- handlers ----------
	userHandler := handlers.NewUserHandler(userClient)
	authHandler := handlers.NewAuthHandler(authClient)

	// ---------- routes ----------
	routes.Setup(app, userHandler, authHandler)

	return app
}
