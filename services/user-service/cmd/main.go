package main

import (
	"log"
	"net"
	grpcclient "wealth-vault/user-service/client"
	config "wealth-vault/user-service/configs"
	userCron "wealth-vault/user-service/internal/delivery/cron"
	handler "wealth-vault/user-service/internal/delivery/grpc"
	"wealth-vault/user-service/internal/event"
	"wealth-vault/user-service/internal/infra"
	repo "wealth-vault/user-service/internal/repository"
	usecase "wealth-vault/user-service/internal/usecase"
	"wealth-vault/user-service/pkg/database"
	userpb "wealth-vault/user-service/pkg/pb/proto/user"
	storageclient "wealth-vault/user-service/pkg/utils"
	mailclient "wealth-vault/user-service/pkg/utils/mail"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfigs()
	database.InitDB(cfg.PostgreSQL)
	db := database.GetDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
	}

	assetClient, err := grpcclient.NewAssetClient(cfg.AssetGRPC.Host, cfg.AssetGRPC.Port)
	if err != nil {
		log.Fatal("user service:", err)
	}

	supabaseClient, err := storageclient.NewStorageClient(cfg.SUPA.URL, cfg.SUPA.Key, cfg.SUPA.Bucket)
	if err != nil {
		log.Fatal("supabase client:", err)
	}

	nc, err := infra.NewNATSConnection(cfg.NATS.Host, cfg.NATS.Port)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	natsPublisher := event.NewPublisher(nc)
	mailClient := mailclient.NewMailClient(cfg.Mail)

	// ------ Repository -------
	urepo := repo.NewUserRepository(db)
	grepo := repo.NewGroupRepository(db)
	irepo := repo.NewShareItemRepository(db)
	mrepo := repo.NewMsgRepository(db)

	// ------ Usecase ------
	guc := usecase.NewGroupUsecase(grepo, urepo, mrepo, supabaseClient, natsPublisher)
	iuc := usecase.NewShareItemUsecase(irepo, grepo, urepo, mrepo, assetClient, mailClient, natsPublisher)
	uuc := usecase.NewUserUsecase(urepo, iuc, supabaseClient, natsPublisher, assetClient)
	muc := usecase.NewMessageUsecase(mrepo, irepo)

	// ------ Handler ------
	uhandler := handler.NewUserGRPCHandler(uuc, guc, iuc, muc)

	// ------ Cronjob ------
	cronjob := userCron.NewAuthCronJob(iuc, uuc)
	cronjob.Start()

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, uhandler)

	log.Printf("🚀 User Service (gRPC) running on :%s", cfg.GRPC.Port)
	log.Fatal(grpcServer.Serve(lis))
}
