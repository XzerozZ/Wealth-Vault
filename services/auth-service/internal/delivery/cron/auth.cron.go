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
	c := cron.New(cron.WithLocation(bkk))
	_, err := c.AddFunc("@daily", j.runCleanup)
	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
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

func (j *AuthCronJob) Stop() {
	j.cron.Stop()
}
