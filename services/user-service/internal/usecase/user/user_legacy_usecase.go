package usecase

import (
	"context"
	"log"
	"wealth-vault/user-service/internal/domain"
	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"
)

func (u *UserUsecase) ProcessLegacyAutoShare(ctx context.Context) error {
	log.Println("⏰ [LEGACY] Start checking for eligible users...")

	users, err := u.userRepo.GetUsersReadyForAutoShare(ctx)
	if err != nil {
		return nil
	}

	log.Printf("🔎 Found %d users eligible for legacy transfer", len(users))
	for _, user := range users {
		u.ProcessSingleUserLegacy(ctx, user)
	}

	return nil
}

func (u *UserUsecase) ProcessSingleUserLegacy(ctx context.Context, user domain.User) error {
	if len(user.Friends) == 0 {
		return nil
	}

	assets, err := u.assetClient.GetAllAssetIDs(ctx, &assetPb.GetMyAssetsRequest{
		UserId: user.ID.String(),
	})
	if err != nil {
		return err
	}

	totalAssets := len(assets.AccountIds) + len(assets.BuildingIds) + len(assets.LandIds) +
		len(assets.CashIds) + len(assets.InsuranceIds) + len(assets.InvestmentIds) + len(assets.LiabilityIds)

	if totalAssets == 0 {
		u.userRepo.MarkAutoShareTriggered(ctx, user.ID)
		return nil
	}

	for _, friend := range user.Friends {
		err := u.itemUC.BatchShareAssets(ctx, domain.BatchShareRequest{
			OwnerID:       user.ID,
			TargetID:      friend.ID,
			AccountIDs:    assets.AccountIds,
			BuildingIDs:   assets.BuildingIds,
			LandIDs:       assets.LandIds,
			CashIDs:       assets.CashIds,
			InsuranceIDs:  assets.InsuranceIds,
			InvestmentIDs: assets.InvestmentIds,
			LiabilityIDs:  assets.LiabilityIds,
		})

		if err != nil {
			log.Printf("⚠️ Failed to share to %s: %v", friend.ID, err)
		}
	}

	u.userRepo.MarkAutoShareTriggered(ctx, user.ID)
	return nil
}
