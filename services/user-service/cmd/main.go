package main

import (
	"log"
	"net"
	config "wealth-vault/user-service/configs"
	handler "wealth-vault/user-service/internal/delivery/grpc"
	repo "wealth-vault/user-service/internal/repository"
	usecase "wealth-vault/user-service/internal/usecase"
	"wealth-vault/user-service/pkg/database"
	userpb "wealth-vault/user-service/pkg/pb/proto/user"
	storageclient "wealth-vault/user-service/pkg/utils"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfigs()
	database.InitDB(cfg.PostgreSQL)
	db := database.GetDB()
	if db == nil {
		log.Fatal("Failed to initialize database")
	}

	supabaseClient, err := storageclient.NewStorageClient(cfg.SUPA.URL, cfg.SUPA.Key, cfg.SUPA.Bucket)
	urepo := repo.NewUserRepository(db)
	grepo := repo.NewGroupRepository(db)
	uuc := usecase.NewUserUsecase(urepo, supabaseClient)
	guc := usecase.NewGroupUsecase(grepo, supabaseClient)
	uhandler := handler.NewUserGRPCHandler(uuc, guc)

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, uhandler)

	log.Printf("🚀 User Service (gRPC) running on :%s", cfg.GRPC.Port)
	log.Fatal(grpcServer.Serve(lis))
}
