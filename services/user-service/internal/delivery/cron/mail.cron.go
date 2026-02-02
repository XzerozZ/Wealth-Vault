package cron

import (
	"context"
	"log"
	"time"
	"wealth-vault/user-service/internal/usecase"

	"github.com/robfig/cron/v3"
)

type MailCronJob struct {
	usecase usecase.ShareItemUsecase
	cron    *cron.Cron
}

func NewAuthCronJob(u usecase.ShareItemUsecase) *MailCronJob {
	return &MailCronJob{
		usecase: u,
		cron:    cron.New(),
	}
}

func (j *MailCronJob) Start() {
	bkk, _ := time.LoadLocation("Asia/Bangkok")
	j.cron = cron.New(cron.WithLocation(bkk))
	_, err := j.cron.AddFunc("0 10 * * *", j.sendEmail)
	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}

	j.cron.Start()
	log.Println("Sending Mail Job started")
}

func (j *MailCronJob) sendEmail() {
	log.Println("Start sending Email")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := j.usecase.ProcessScheduledEmails(ctx); err != nil {
		log.Printf("Sending Email failed : %v\n", err)
	} else {
		log.Printf("Sending completed")
	}
}

func (j *MailCronJob) Stop() {
	j.cron.Stop()
}
