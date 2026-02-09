package main

import (
	"log"
	"net"
	grpcclient "wealth-vault/asset-service/client"
	"wealth-vault/asset-service/configs"
	assetCron "wealth-vault/asset-service/internal/delivery/cron"
	handler "wealth-vault/asset-service/internal/delivery/grpc"
	"wealth-vault/asset-service/internal/event"
	"wealth-vault/asset-service/internal/infra"
	repo "wealth-vault/asset-service/internal/repository"
	usecase "wealth-vault/asset-service/internal/usecase"
	"wealth-vault/asset-service/pkg/database"
	assetpb "wealth-vault/asset-service/pkg/pb/proto/asset"
	storageclient "wealth-vault/asset-service/pkg/utils"

	"google.golang.org/grpc"
)

func main() {
	cfg := configs.LoadConfigs()
	database.InitDB(cfg.PostgreSQL)
	db := database.GetDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
	}

	nc, err := infra.NewNATSConnection(cfg.NATS.Host, cfg.NATS.Port)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	supabaseClient, err := storageclient.NewStorageClient(cfg.SUPA.URL, cfg.SUPA.Key, cfg.SUPA.Bucket)
	if err != nil {
		log.Fatal("supabase client:", err)
	}

	userClient, err := grpcclient.NewUserClient(cfg.UserGRPC.Host, cfg.UserGRPC.Port)
	if err != nil {
		log.Fatal("user service:", err)
	}

	natsPublisher := event.NewPublisher(nc)
	// ------ Repository ------
	assRepo := repo.NewAssetRepository(db)
	accRepo := repo.NewAccountRepository(db)
	cashRepo := repo.NewCashRepository(db)
	inRepo := repo.NewInvestmentRepository(db)
	buRepo := repo.NewBuildingRepository(db)
	landRepo := repo.NewLandRepository(db)
	insRepo := repo.NewInsuranceRepository(db)
	fileRepo := repo.NewFileRepository(db)
	liaRepo := repo.NewLiabilityRepository(db)

	// ------ Usecase ------
	assUC := usecase.NewAssetUsecase(accRepo, buRepo, cashRepo, insRepo, inRepo, landRepo, liaRepo, assRepo)
	accUC := usecase.NewAccountUsecase(accRepo, fileRepo, supabaseClient, userClient)
	cashUC := usecase.NewCashUsecase(cashRepo, fileRepo, supabaseClient, userClient)
	inUC := usecase.NewInvestmentUsecase(inRepo, fileRepo, supabaseClient, userClient)
	buUC := usecase.NewBuildingUsecase(buRepo, fileRepo, supabaseClient, userClient)
	landUC := usecase.NewLandUsecase(landRepo, fileRepo, supabaseClient, userClient)
	insUC := usecase.NewInsuranceUsecase(insRepo, fileRepo, supabaseClient, natsPublisher, userClient)
	liaUC := usecase.NewLiabilityUsecase(liaRepo, fileRepo, supabaseClient, userClient)

	// ------ Handler ------
	assetHandler := handler.NewAssetGRPCHandler(assUC, accUC, cashUC, inUC, buUC, landUC, insUC, liaUC)

	// ----- Cronjob ------
	cronjob := assetCron.NewAssetCronJob(accUC, buUC, cashUC, insUC, inUC, landUC, liaUC)
	cronjob.Start()

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	assetpb.RegisterAssetServiceServer(grpcServer, assetHandler)

	log.Printf("🚀 Asset Service (gRPC) running on :%s", cfg.GRPC.Port)
	log.Fatal(grpcServer.Serve(lis))
}
