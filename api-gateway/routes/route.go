package routes

import (
	"wealth-vault/api-gateway/configs"
	"wealth-vault/api-gateway/handlers"
	"wealth-vault/api-gateway/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Setup(
	app *fiber.App,
	jwt configs.JWT,
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	accHandler *handlers.AccountHandler,
	cashHandler *handlers.CashHandler,
	inHandler *handlers.InvestmentHandler,
	buHandler *handlers.BuildingHandler,
	landHandler *handlers.LandHandler,
	insHandler *handlers.InsuranceHandler,
	liaHandler *handlers.LiabilityHandler,
	groupHandler *handlers.GroupHandler,
	groupitemHandler *handlers.GroupItemHandler,
	notificationHandler *handlers.NotificationHandler,
	msgHandler *handlers.MessageHandler,
	infoHandler *handlers.InfoHandler,
) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.RegisterLocal)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/forgot/password", authHandler.ForgotPassword)
	auth.Post("/forgot/otp", authHandler.VerifyForgotPasswordOTP)
	auth.Patch("/reset/password", authHandler.ResetPassword)

	api.Get("/user", middleware.JWTMiddleware(jwt), userHandler.GetUser)
	api.Patch("/user", middleware.JWTMiddleware(jwt), userHandler.UpdateUser)
	api.Post("/friend", middleware.JWTMiddleware(jwt), userHandler.AddFriend)
	api.Post("/friend/accept", middleware.JWTMiddleware(jwt), userHandler.AcceptFriend)
	api.Get("/friend", middleware.JWTMiddleware(jwt), userHandler.GetFriendList)
	api.Get("/friend/:id/profile", middleware.JWTMiddleware(jwt), userHandler.GetFriendProfile)
	api.Get("/friend/pending", middleware.JWTMiddleware(jwt), userHandler.GetPendingRequests)
	api.Post("/closefriend", middleware.JWTMiddleware(jwt), userHandler.ToggleCloseFriend)
	api.Get("/closefriend", middleware.JWTMiddleware(jwt), userHandler.GetCloseFriends)
	api.Get("/friend/:id/item", middleware.JWTMiddleware(jwt), groupitemHandler.GetFriendItems)
	api.Post("/group/:id/addmember", middleware.JWTMiddleware(jwt), groupHandler.AddMember)
	api.Post("/group/:id/grantaccess", middleware.JWTMiddleware(jwt), groupHandler.GrantAccess)
	api.Delete("/group/:id/removemember", middleware.JWTMiddleware(jwt), groupHandler.RemoveMember)
	api.Delete("/group/:id/leave", middleware.JWTMiddleware(jwt), groupHandler.LeaveGroup)

	api.Get("/dashboard", middleware.JWTMiddleware(jwt), infoHandler.Dashboard)

	api.Post("/group", middleware.JWTMiddleware(jwt), groupHandler.CreateGroup)
	api.Get("/group", middleware.JWTMiddleware(jwt), groupHandler.AllGetGroup)
	api.Get("/group/detail/:id", middleware.JWTMiddleware(jwt), groupHandler.GetGroup)
	api.Get("/group/member/:id", middleware.JWTMiddleware(jwt), groupHandler.GetMember)
	api.Patch("/group/:id", middleware.JWTMiddleware(jwt), groupHandler.UpdateGroup)
	api.Delete("/group/:id", middleware.JWTMiddleware(jwt), groupHandler.DeleteGroup)
	api.Post("/share/item", middleware.JWTMiddleware(jwt), groupitemHandler.ShareItem)
	api.Get("/group/:id/item", middleware.JWTMiddleware(jwt), groupitemHandler.GetGroupItems)
	api.Delete("/group/item/:id", middleware.JWTMiddleware(jwt), groupitemHandler.UnsharedItem)
	api.Delete("/friend/item/:id", middleware.JWTMiddleware(jwt), groupitemHandler.UnsharedIteminFriend)

	api.Post("/asset/account", middleware.JWTMiddleware(jwt), accHandler.CreateAccount)
	api.Get("/asset/account", middleware.JWTMiddleware(jwt), accHandler.GetAccount)
	api.Get("/asset/account/:id", middleware.JWTMiddleware(jwt), accHandler.GetAccountByID)
	api.Patch("/asset/account/:id", middleware.JWTMiddleware(jwt), accHandler.UpdateAccount)
	api.Delete("/asset/account/:id", middleware.JWTMiddleware(jwt), accHandler.DeleteAccount)

	api.Post("/asset/invest", middleware.JWTMiddleware(jwt), inHandler.CreateInvestment)
	api.Get("/asset/invest", middleware.JWTMiddleware(jwt), inHandler.GetInvestment)
	api.Get("/asset/invest/:id", middleware.JWTMiddleware(jwt), inHandler.GetInvestmentByID)
	api.Patch("/asset/invest/:id", middleware.JWTMiddleware(jwt), inHandler.UpdateInvestment)
	api.Delete("/asset/invest/:id", middleware.JWTMiddleware(jwt), inHandler.DeleteInvestment)

	api.Post("/asset/cash", middleware.JWTMiddleware(jwt), cashHandler.CreateCash)
	api.Get("/asset/cash", middleware.JWTMiddleware(jwt), cashHandler.GetCash)
	api.Get("/asset/cash/:id", middleware.JWTMiddleware(jwt), cashHandler.GetCashByID)
	api.Patch("/asset/cash/:id", middleware.JWTMiddleware(jwt), cashHandler.UpdateCash)
	api.Delete("/asset/cash/:id", middleware.JWTMiddleware(jwt), cashHandler.DeleteCash)

	api.Post("/asset/building", middleware.JWTMiddleware(jwt), buHandler.CreateBuilding)
	api.Get("/asset/building", middleware.JWTMiddleware(jwt), buHandler.GetBuilding)
	api.Get("/asset/building/:id", middleware.JWTMiddleware(jwt), buHandler.GetBuildingByID)
	api.Patch("/asset/building/:id", middleware.JWTMiddleware(jwt), buHandler.UpdateBuilding)
	api.Delete("/asset/building/:id", middleware.JWTMiddleware(jwt), buHandler.DeleteBuilding)

	api.Post("/asset/land", middleware.JWTMiddleware(jwt), landHandler.CreateLand)
	api.Get("/asset/land", middleware.JWTMiddleware(jwt), landHandler.GetLand)
	api.Get("/asset/land/:id", middleware.JWTMiddleware(jwt), landHandler.GetLandByID)
	api.Patch("/asset/land/:id", middleware.JWTMiddleware(jwt), landHandler.UpdateLand)
	api.Delete("/asset/land/:id", middleware.JWTMiddleware(jwt), landHandler.DeleteLand)

	api.Post("/asset/insurance", middleware.JWTMiddleware(jwt), insHandler.CreateInsurance)
	api.Get("/asset/insurance", middleware.JWTMiddleware(jwt), insHandler.GetInsurance)
	api.Get("/asset/insurance/:id", middleware.JWTMiddleware(jwt), insHandler.GetInsuranceByID)
	api.Patch("/asset/insurance/:id", middleware.JWTMiddleware(jwt), insHandler.UpdateInsurance)
	api.Delete("/asset/insurance/:id", middleware.JWTMiddleware(jwt), insHandler.DeleteInsurance)

	api.Post("/lia", middleware.JWTMiddleware(jwt), liaHandler.CreateLiability)
	api.Get("/lia", middleware.JWTMiddleware(jwt), liaHandler.GetLiability)
	api.Get("/lia/:id", middleware.JWTMiddleware(jwt), liaHandler.GetLiabilityByID)
	api.Patch("/lia/:id", middleware.JWTMiddleware(jwt), liaHandler.UpdateLiability)
	api.Delete("/lia/:id", middleware.JWTMiddleware(jwt), liaHandler.DeleteLiability)

	api.Get("/share/item/:type/:id/shared-targets", middleware.JWTMiddleware(jwt), groupitemHandler.GetItemSharedTargets)
	api.Get("share/:type/:id/selection", middleware.JWTMiddleware(jwt), groupitemHandler.GetItemsForSelection)

	api.Get("/ws", middleware.TokenFromQuery, middleware.JWTMiddleware(jwt), websocket.New(notificationHandler.ProxyWebSocket))
	api.All("/notifications/*", middleware.JWTMiddleware(jwt), notificationHandler.ProxyAPI)

	api.Get("/group/:id/msg", middleware.JWTMiddleware(jwt), msgHandler.GetGroupMessages)
	api.Get("/friend/:id/msg", middleware.JWTMiddleware(jwt), msgHandler.GetPrivateMessages)
}
