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
	liaHandler *handlers.LiabilityHandler,
	groupHandler *handlers.GroupHandler,
) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.RegisterLocal)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	api.Get("/user", middleware.JWTMiddleware(jwt), userHandler.GetUser)
	api.Patch("/user", middleware.JWTMiddleware(jwt), userHandler.UpdateUser)
	api.Post("/friend/:friendID", middleware.JWTMiddleware(jwt), userHandler.AddFriend)
	api.Get("/friend", middleware.JWTMiddleware(jwt), userHandler.GetFriendList)

	api.Post("/group", middleware.JWTMiddleware(jwt), groupHandler.CreateGroup)
	api.Get("/group/detail/:id", middleware.JWTMiddleware(jwt), groupHandler.GetGroup)
	api.Get("/group/member/:id", middleware.JWTMiddleware(jwt), groupHandler.GetMember)
	api.Patch("/group/:id", middleware.JWTMiddleware(jwt), groupHandler.UpdateGroup)

	api.Post("/asset", middleware.JWTMiddleware(jwt), assetHandler.CreateAsset)
	api.Get("/asset", middleware.JWTMiddleware(jwt), assetHandler.GetAsset)
	api.Get("/asset/:id", middleware.JWTMiddleware(jwt), assetHandler.GetAssetByID)
	api.Patch("/asset/:id", middleware.JWTMiddleware(jwt), assetHandler.UpdateAsset)
	api.Delete("/asset/:id", middleware.JWTMiddleware(jwt), assetHandler.DeleteAsset)

	api.Post("/lia", middleware.JWTMiddleware(jwt), liaHandler.CreateLiability)
	api.Get("/lia", middleware.JWTMiddleware(jwt), liaHandler.GetLiability)
	api.Get("/lia/:id", middleware.JWTMiddleware(jwt), liaHandler.GetLiabilityByID)
	api.Patch("/lia/:id", middleware.JWTMiddleware(jwt), liaHandler.UpdateLiability)
	api.Delete("/lia/:id", middleware.JWTMiddleware(jwt), liaHandler.DeleteLiability)
}
