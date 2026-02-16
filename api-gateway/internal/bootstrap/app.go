package bootstrap

import (
	"log"
	"time"

	grpcclient "wealth-vault/api-gateway/client"
	"wealth-vault/api-gateway/configs"
	"wealth-vault/api-gateway/handlers"
	"wealth-vault/api-gateway/internal/middleware"
	storageclient "wealth-vault/api-gateway/pkg/utils"
	"wealth-vault/api-gateway/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func InitApp(cfg *configs.Configs) *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit:         2 * 1024 * 1024 * 1024,
		ReadTimeout:       10 * time.Minute,
		WriteTimeout:      30 * time.Minute,
		IdleTimeout:       10 * time.Minute,
		StreamRequestBody: true,
	})

	app.Use(recover.New())
	app.Use(middleware.Cors())
	supabaseClient, err := storageclient.NewStorageClient(cfg.SUPA.URL, cfg.SUPA.Key, cfg.SUPA.Bucket)
	// ---------- gRPC clients ----------
	userClient, err := grpcclient.NewUserClient(cfg.UserService.Host, cfg.UserService.Port)
	if err != nil {
		log.Fatal("user service:", err)
	}

	authClient, err := grpcclient.NewAuthClient(cfg.AuthService.Host, cfg.AuthService.Port)
	if err != nil {
		log.Fatal("auth service:", err)
	}

	assetClient, err := grpcclient.NewAssetClient(cfg.AssetService.Host, cfg.AssetService.Port)
	if err != nil {
		log.Fatal("asset service:", err)
	}

	// ---------- handlers ----------
	userHandler := handlers.NewUserHandler(userClient, supabaseClient)
	authHandler := handlers.NewAuthHandler(authClient)
	accHandler := handlers.NewAccountHandler(assetClient, supabaseClient)
	cashHandler := handlers.NewCashHandler(assetClient, supabaseClient)
	inHandler := handlers.NewInvestmentHandler(assetClient, supabaseClient)
	buHandler := handlers.NewBuildingHandler(assetClient, supabaseClient)
	landHandler := handlers.NewLandHandler(assetClient, supabaseClient)
	insHandler := handlers.NewInsuranceHandler(assetClient, supabaseClient)
	liaHandler := handlers.NewLiabilityHandler(assetClient, supabaseClient)
	groupHandler := handlers.NewGroupHandler(userClient, supabaseClient)
	groupItemHandler := handlers.NewGroupItemHandler(userClient, assetClient)
	notificationHandler := handlers.NewNotificationHandler(cfg)
	msgHandler := handlers.NewMessageHandlerr(userClient)
	infoHandler := handlers.NewInfoHandler(assetClient)

	// ---------- routes ----------
	routes.Setup(app,
		cfg.JWT,
		userHandler,
		authHandler,
		accHandler,
		cashHandler,
		inHandler,
		buHandler,
		landHandler,
		insHandler,
		liaHandler,
		groupHandler,
		groupItemHandler,
		notificationHandler,
		msgHandler,
		infoHandler,
	)

	return app
}
