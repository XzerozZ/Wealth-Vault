package routes

import (
	"wealth-vault/api-gateway/handlers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App,
	userHandler *handlers.UserHandler,
) {
	api := app.Group("/api")

	api.Post("/user", userHandler.CreateUser)
}
