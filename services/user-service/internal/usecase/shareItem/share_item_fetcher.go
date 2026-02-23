package usecase

import (
	"context"
	"log"
	"strings"
	"sync"
	"wealth-vault/user-service/internal/domain"
	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils/mapper"
)

func (u *ShareItemUsecase) FetchAssetPreviews(ctx context.Context, items []domain.SharedItemSummary) (map[string]*pb.AssetPreview, error) {
	idsByType := make(map[string][]string)
	for _, item := range items {
		key := strings.ToLower(item.EntityType)
		idsByType[key] = append(idsByType[key], item.EntityID)
	}

	previewMap := make(map[string]*pb.AssetPreview)
	var wg sync.WaitGroup
	var mu sync.Mutex

	runFetch := func(assetType string, fetcher func() error) {
		if len(idsByType[assetType]) == 0 {
			return
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fetcher(); err != nil {
				log.Printf("⚠️ Failed to fetch %s assets: %v", assetType, err)
			}
		}()
	}

	reqFor := func(t string) *assetPb.GetBatchIdsRequest {
		return &assetPb.GetBatchIdsRequest{Ids: idsByType[t]}
	}

	runFetch(AssetTypeBuilding, func() error {
		res, err := u.assetClient.GetBatchBuilding(ctx, reqFor(AssetTypeBuilding))
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		for _, b := range res.Building {
			previewMap[b.Id] = mapper.MapBuildingToPreview(b)
		}
		return nil
	})

	runFetch(AssetTypeLand, func() error {
		res, err := u.assetClient.GetBatchLand(ctx, reqFor(AssetTypeLand))
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		for _, l := range res.Land {
			previewMap[l.Id] = mapper.MapLandToPreview(l)
		}
		return nil
	})

	runFetch(AssetTypeAccount, func() error {
		res, err := u.assetClient.GetBatchAccount(ctx, reqFor(AssetTypeAccount))
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		for _, a := range res.Account {
			previewMap[a.Id] = mapper.MapAccountToPreview(a)
		}
		return nil
	})

	runFetch(AssetTypeCash, func() error {
		res, err := u.assetClient.GetBatchCash(ctx, reqFor(AssetTypeCash))
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		for _, c := range res.Cash {
			previewMap[c.Id] = mapper.MapCashToPreview(c)
		}
		return nil
	})

	runFetch(AssetTypeInsurance, func() error {
		res, err := u.assetClient.GetBatchInsurance(ctx, reqFor(AssetTypeInsurance))
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		for _, i := range res.Insurance {
			previewMap[i.Id] = mapper.MapInsuranceToPreview(i)
		}
		return nil
	})

	runFetch(AssetTypeInvestment, func() error {
		res, err := u.assetClient.GetBatchInvestment(ctx, reqFor(AssetTypeInvestment))
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		for _, inv := range res.Invest {
			previewMap[inv.Id] = mapper.MapInvestmentToPreview(inv)
		}
		return nil
	})

	runFetch(AssetTypeLiability, func() error {
		res, err := u.assetClient.GetBatchLiability(ctx, reqFor(AssetTypeLiability))
		if err != nil {
			return err
		}
		mu.Lock()
		defer mu.Unlock()
		for _, l := range res.Liability {
			previewMap[l.Id] = mapper.MapLiabilityToPreview(l)
		}
		return nil
	})

	wg.Wait()
	return previewMap, nil
}
