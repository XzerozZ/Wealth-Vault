package cron

import (
	"context"
	"log"
	"time"
	"wealth-vault/auth-service/internal/usecase"

	"github.com/robfig/cron/v3"
)

type AuthCronJob struct {
	usecase usecase.AuthUsecase
	cron    *cron.Cron
}

func NewAuthCronJob(u usecase.AuthUsecase) *AuthCronJob {
	return &AuthCronJob{
		usecase: u,
		cron:    cron.New(),
	}
}

func (j *AuthCronJob) Start() {
	bkk, _ := time.LoadLocation("Asia/Bangkok")
	j.cron = cron.New(cron.WithLocation(bkk))
	if _, err := j.cron.AddFunc("@daily", j.runCleanup); err != nil {
		log.Fatalf("Error adding cleanup job: %v", err)
	}

	if _, err := j.cron.AddFunc("@every 1h", j.runCleanupOTP); err != nil {
		log.Fatalf("Error adding OTP cleanup job: %v", err)
	}

	j.cron.Start()
	log.Println("Auth Cleanup Job started")
}

func (j *AuthCronJob) runCleanup() {
	log.Println("Start session cleanup")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := j.usecase.CleanupSessions(ctx); err != nil {
		log.Printf("Cleanup failed: %v\n", err)
	} else {
		log.Println("Cleanup completed successfully")
	}
}

func (j *AuthCronJob) runCleanupOTP() {
	log.Println("Start session cleanup OTP")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := j.usecase.DeleteExpiredOTPs(ctx); err != nil {
		log.Printf("Cleanup failed: %v\n", err)
	} else {
		log.Println("Cleanup completed successfully")
	}
}

func (j *AuthCronJob) Stop() {
	j.cron.Stop()
}
