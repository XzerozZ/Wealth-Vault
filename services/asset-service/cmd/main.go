package main

import (
	"log"
	"net"
	"wealth-vault/asset-service/configs"
	handler "wealth-vault/asset-service/internal/delivery/grpc"
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

	supabaseClient, err := storageclient.NewStorageClient(cfg.SUPA.URL, cfg.SUPA.Key, cfg.SUPA.Bucket)
	// ------ Repository------
	assRepo := repo.NewAssetRepository(db)
	accRepo := repo.NewAccountRepository(db)
	cashRepo := repo.NewCashRepository(db)
	inRepo := repo.NewInvestmentRepository(db)
	buRepo := repo.NewBuildingRepository(db)
	landRepo := repo.NewLandRepository(db)
	insRepo := repo.NewInsuranceRepository(db)
	fileRepo := repo.NewFileRepository(db)
	liaRepo := repo.NewLiabilityRepository(db)

	// ------ Usecase------
	assUC := usecase.NewAssetUsecase(assRepo)
	accUC := usecase.NewAccountUsecase(accRepo, fileRepo, supabaseClient)
	cashUC := usecase.NewCashUsecase(cashRepo, fileRepo, supabaseClient)
	inUC := usecase.NewInvestmentUsecase(inRepo, fileRepo, supabaseClient)
	buUC := usecase.NewBuildingUsecase(buRepo, fileRepo, supabaseClient)
	landUC := usecase.NewLandUsecase(landRepo, fileRepo, supabaseClient)
	insUC := usecase.NewInsuranceUsecase(insRepo, fileRepo, supabaseClient)
	liaUC := usecase.NewLiabilityUsecase(liaRepo, fileRepo, supabaseClient)

	// ------Handler------
	assetHandler := handler.NewAssetGRPCHandler(assUC, accUC, cashUC, inUC, buUC, landUC, insUC, liaUC)

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	assetpb.RegisterAssetServiceServer(grpcServer, assetHandler)

	log.Printf("🚀 Asset Service (gRPC) running on :%s", cfg.GRPC.Port)
	log.Fatal(grpcServer.Serve(lis))
}
