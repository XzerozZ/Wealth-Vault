package cron

import (
	"context"
	"log"
	"time"
	"wealth-vault/user-service/internal/usecase"

	"github.com/robfig/cron/v3"
)

type MailCronJob struct {
	itemusecase usecase.ShareItemUsecase
	userusecase usecase.UserUsecase
	cron        *cron.Cron
}

func NewAuthCronJob(u usecase.ShareItemUsecase, uu usecase.UserUsecase) *MailCronJob {
	return &MailCronJob{
		itemusecase: u,
		userusecase: uu,
		cron:        cron.New(),
	}
}

func (j *MailCronJob) Start() {
	bkk, _ := time.LoadLocation("Asia/Bangkok")
	j.cron = cron.New(cron.WithLocation(bkk))
	_, err := j.cron.AddFunc("0 10 * * *", j.sendEmail)
	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}

	_, err = j.cron.AddFunc("* 10 * * *", j.AutoShareTrigger)
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

	if err := j.itemusecase.ProcessScheduledEmails(ctx); err != nil {
		log.Printf("Sending Email failed : %v\n", err)
	} else {
		log.Printf("Sending completed")
	}
}

func (j *MailCronJob) AutoShareTrigger() {
	log.Println("Start Auto Share Trigger")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := j.userusecase.ProcessLegacyAutoShare(ctx); err != nil {
		log.Printf("Auto Share failed : %v\n", err)
	} else {
		log.Printf("Auto Share completed")
	}
}

func (j *MailCronJob) Stop() {
	j.cron.Stop()
}
