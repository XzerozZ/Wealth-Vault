package main

import (
	"log"
	"net"
	"wealth-vault/asset-service/configs"
	assetHandler "wealth-vault/asset-service/internal/delivery/grpc"
	assetRepo "wealth-vault/asset-service/internal/repository"
	assetUsecase "wealth-vault/asset-service/internal/usecase"
	"wealth-vault/asset-service/pkg/database"
	assetpb "wealth-vault/asset-service/pkg/pb/proto/asset"

	"google.golang.org/grpc"
)

func main() {
	cfg := configs.LoadConfigs()
	database.InitDB(cfg.PostgreSQL)
	db := database.GetDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
	}

	repo := assetRepo.NewAssetRepository(db)
	uc := assetUsecase.NewAssetUsecase(repo)
	handler := assetHandler.NewAssetGRPCHandler(uc)

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	assetpb.RegisterAssetServiceServer(grpcServer, handler)

	log.Printf("🚀 Asset Service (gRPC) running on :%s", cfg.GRPC.Port)
	log.Fatal(grpcServer.Serve(lis))
}
