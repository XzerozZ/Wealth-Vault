package main

import (
	"log"
	"net"
	grpcclient "wealth-vault/auth-service/client"
	"wealth-vault/auth-service/configs"
	authHandler "wealth-vault/auth-service/internal/handler/grpc"
	authRepo "wealth-vault/auth-service/internal/repository"
	authUsecase "wealth-vault/auth-service/internal/usecase"
	"wealth-vault/auth-service/pkg/database"
	authpb "wealth-vault/auth-service/pkg/pb/proto/auth"

	"google.golang.org/grpc"
)

func main() {
	cfg := configs.LoadConfigs()
	database.InitDB(cfg.PostgreSQL)
	db := database.GetDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
	}

	userClient, err := grpcclient.NewUserClient(cfg.UserGRPC.Host, cfg.UserGRPC.Port)
	if err != nil {
		log.Fatal("user service:", err)
	}

	repo := authRepo.NewAuthRepository(db)
	uc := authUsecase.NewAuthUsecase(repo, userClient)
	handler := authHandler.NewAuthGRPCHandler(uc)

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, handler)

	log.Printf("🚀 Auth Service (gRPC) running on :%s", cfg.GRPC.Port)
	log.Fatal(grpcServer.Serve(lis))
}
