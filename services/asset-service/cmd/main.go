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
	assetRepo := repo.NewAssetRepository(db)
	fileRepo := repo.NewFileRepository(db)
	// ------ Usecase------
	uc := usecase.NewAssetUsecase(assetRepo, fileRepo, supabaseClient)
	// ------Handler------
	handler := handler.NewAssetGRPCHandler(uc)

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	assetpb.RegisterAssetServiceServer(grpcServer, handler)

	log.Printf("🚀 Asset Service (gRPC) running on :%s", cfg.GRPC.Port)
	log.Fatal(grpcServer.Serve(lis))
}
