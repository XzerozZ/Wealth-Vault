package main

import (
	"log"
	"net"
	config "wealth-vault/user-service/configs"
	userHandler "wealth-vault/user-service/internal/delivery/grpc"
	userRepo "wealth-vault/user-service/internal/repository"
	userUsecase "wealth-vault/user-service/internal/usecase"
	"wealth-vault/user-service/pkg/database"
	userpb "wealth-vault/user-service/pkg/pb/proto/user"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfigs()
	database.InitDB(cfg.PostgreSQL)
	db := database.GetDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
	}

	repo := userRepo.NewUserRepository(db)
	uc := userUsecase.NewUserUsecase(repo)
	handler := userHandler.NewUserGRPCHandler(uc)

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, handler)

	log.Printf("🚀 User Service (gRPC) running on :%s", cfg.GRPC.Port)
	log.Fatal(grpcServer.Serve(lis))
}
