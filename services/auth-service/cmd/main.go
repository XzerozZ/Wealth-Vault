package main

import (
	"log"
	"net"
	"os"
	"time"
	grpcclient "wealth-vault/auth-service/client"
	"wealth-vault/auth-service/configs"
	authCron "wealth-vault/auth-service/internal/delivery/cron"
	authHandler "wealth-vault/auth-service/internal/delivery/grpc"
	"wealth-vault/auth-service/internal/infra/database"
	authRepo "wealth-vault/auth-service/internal/repository"
	authUsecase "wealth-vault/auth-service/internal/usecase"
	"wealth-vault/auth-service/pkg/google"
	mailclient "wealth-vault/auth-service/pkg/mail"
	authpb "wealth-vault/auth-service/pkg/pb/proto/auth"
	authToken "wealth-vault/auth-service/pkg/token"

	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--healthcheck" {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+os.Getenv("GRPC_PORT"), time.Second)
		if err != nil {
			os.Exit(1)
		}
		_ = conn.Close()
		return
	}
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

	mailClient, err := mailclient.NewMailClient(cfg.Mail)
	if err != nil {
		log.Fatal("mail client failed:", err)
	}

	googleVal := google.NewGoogleValidator(cfg.GoogleClient.URL)

	repo := authRepo.NewAuthRepository(db)
	token := authToken.NewJWTGenerate(cfg.JWT.Secret)
	uc := authUsecase.NewAuthUsecase(repo, userClient, token, mailClient, googleVal)
	uc.ConfigureLINE(os.Getenv("LINE_CHANNEL_ID"), authUsecase.NewLINEIDTokenVerifier())
	handler := authHandler.NewAuthGRPCHandler(uc)
	cronJob := authCron.NewAuthCronJob(uc)

	cronJob.Start()

	lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, handler)

	log.Printf("🚀 Auth Service (gRPC) running on :%s", cfg.GRPC.Port)
	log.Fatal(grpcServer.Serve(lis))
}
