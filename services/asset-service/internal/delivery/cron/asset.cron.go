package cron

import (
	"context"
	"log"
	"time"
	usecase "wealth-vault/asset-service/internal/usecase/interface"

	"github.com/robfig/cron/v3"
)

type AssetCronJob struct {
	a    usecase.AccountUsecase
	b    usecase.BuildingUsecase
	c    usecase.CashUsecase
	ins  usecase.InsuranceUsecase
	inv  usecase.InvestmentUsecase
	la   usecase.LandUsecase
	lia  usecase.LiabilityUsecase
	cron *cron.Cron
}

func NewAssetCronJob(
	aUC usecase.AccountUsecase,
	bUC usecase.BuildingUsecase,
	cUC usecase.CashUsecase,
	insUC usecase.InsuranceUsecase,
	invUC usecase.InvestmentUsecase,
	laUC usecase.LandUsecase,
	liaUC usecase.LiabilityUsecase,
) *AssetCronJob {
	return &AssetCronJob{
		a:    aUC,
		b:    bUC,
		c:    cUC,
		ins:  insUC,
		inv:  invUC,
		la:   laUC,
		lia:  liaUC,
		cron: cron.New(),
	}
}

func (j *AssetCronJob) Start() {
	bkk, _ := time.LoadLocation("Asia/Bangkok")
	j.cron = cron.New(cron.WithLocation(bkk))
	_, err := j.cron.AddFunc("* 10 * * *", j.sendNoti)
	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}

	_, err = j.cron.AddFunc("0 3 * * *", j.CleanUpExpiredAsset)
	if err != nil {
		log.Fatalf("Error adding cleanup cron job: %v", err)
	}

	j.cron.Start()
	log.Println("Sending Insurance Expiring Notification Job started")
}

func (j *AssetCronJob) CleanUpExpiredAsset() {
	log.Println("Starting Daily Cleanup Process...")
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	j.safeExecute("Account", func() {
		j.a.CleanupExpiredAccounts(ctx)
	})

	j.safeExecute("Building", func() {
		j.b.CleanupExpiredBuildings(ctx)
	})

	j.safeExecute("Land", func() {
		j.la.CleanupExpiredLands(ctx)
	})

	j.safeExecute("Cash", func() {
		j.c.CleanupExpiredCashes(ctx)
	})

	j.safeExecute("Insurance", func() {
		j.ins.CleanupExpiredInsurances(ctx)
	})

	j.safeExecute("Investment", func() {
		j.inv.CleanupExpiredInvestment(ctx)
	})

	j.safeExecute("Liability", func() {
		j.lia.CleanupExpiredLiabilities(ctx)
	})

	duration := time.Since(startTime)
	log.Printf("Cleanup Process Finished in %v", duration)
}

func (j *AssetCronJob) sendNoti() {
	log.Println("Start Sending Insurance Expiring Notification Job")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := j.ins.CheckExpiringInsurances(ctx); err != nil {
		log.Printf("Sending Noti failed : %v\n", err)
	} else {
		log.Printf("Sending completed")
	}
}

func (j *AssetCronJob) safeExecute(assetName string, job func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf(" Error while cleaning up %s: %v", assetName, r)
		}
	}()

	log.Printf("⏳ Cleaning up %s...", assetName)
	job()
}

func (j *AssetCronJob) Stop() {
	j.cron.Stop()
}
