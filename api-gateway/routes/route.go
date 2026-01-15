package routes

import (
	"wealth-vault/api-gateway/configs"
	"wealth-vault/api-gateway/handlers"
	"wealth-vault/api-gateway/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(
	app *fiber.App,
	jwt configs.JWT,
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	assetHandler *handlers.AssetHandler,
) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.RegisterLocal)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	api.Get("/user", middleware.JWTMiddleware(jwt), userHandler.GetUser)
	api.Patch("/user", middleware.JWTMiddleware(jwt), userHandler.UpdateUser)

	api.Post("/cash", middleware.JWTMiddleware(jwt), assetHandler.CreateCash)
	api.Get("/cash", middleware.JWTMiddleware(jwt), assetHandler.GetCash)
	api.Get("/cash/:id", middleware.JWTMiddleware(jwt), assetHandler.GetCashByID)
	api.Patch("/cash/:id", middleware.JWTMiddleware(jwt), assetHandler.UpdateCash)
	api.Delete("/cash/:id", middleware.JWTMiddleware(jwt), assetHandler.DeleteCash)
}
