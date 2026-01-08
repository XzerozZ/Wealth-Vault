package routes

import (
	"wealth-vault/api-gateway/handlers"

	"github.com/gofiber/fiber/v2"
)

func Setup(
	app *fiber.App,
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.RegisterLocal)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	api.Post("/user", userHandler.CreateUser)
}
