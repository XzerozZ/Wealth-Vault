package usecase

import (
	"context"
	repo "wealth-vault/asset-service/internal/repository/interface"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils"
	"wealth-vault/asset-service/pkg/utils/mapper"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type AssetUsecase struct {
	a   repo.AccountRepository
	b   repo.BuildingRepository
	c   repo.CashRepository
	ins repo.InsuranceRepository
	inv repo.InvestmentRepository
	la  repo.LandRepository
	lia repo.LiabilityRepository
	r   repo.AssetRepository
}

func NewAssetUsecase(
	a repo.AccountRepository,
	b repo.BuildingRepository,
	c repo.CashRepository,
	ins repo.InsuranceRepository,
	inv repo.InvestmentRepository,
	la repo.LandRepository,
	lia repo.LiabilityRepository,
	r repo.AssetRepository,
) AssetUsecase {
	return AssetUsecase{
		a:   a,
		b:   b,
		c:   c,
		ins: ins,
		inv: inv,
		la:  la,
		lia: lia,
		r:   r,
	}
}

func (u *AssetUsecase) CheckExists(ctx context.Context, req *pb.CheckAssetRequest) (*pb.CheckAssetResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	res, err := u.r.CheckExists(ctx, req.Type, id, uid)
	if err != nil {
		return nil, err
	}

	return &pb.CheckAssetResponse{
		Exists: res,
	}, nil
}

func (u *AssetUsecase) GetAllAssetIDs(ctx context.Context, req *pb.GetMyAssetsRequest) (*pb.GetMyAssetsResponse, error) {
	uid, _ := uuid.Parse(req.UserId)
	resp := &pb.GetMyAssetsResponse{}
	var g errgroup.Group
	g.Go(func() error {
		items, err := u.a.GetAccountByUserID(ctx, uid)
		if err != nil {
			return err
		}
		for _, item := range items {
			resp.AccountIds = append(resp.AccountIds, item.ID.String())
		}
		return nil
	})

	g.Go(func() error {
		items, err := u.b.GetBuildingByUserID(ctx, uid)
		if err != nil {
			return err
		}
		for _, item := range items {
			resp.BuildingIds = append(resp.BuildingIds, item.ID.String())
		}
		return nil
	})

	g.Go(func() error {
		items, err := u.c.GetCashByUserID(ctx, uid)
		if err != nil {
			return err
		}
		for _, item := range items {
			resp.CashIds = append(resp.CashIds, item.ID.String())
		}
		return nil
	})

	g.Go(func() error {
		items, err := u.ins.GetInsuranceByUserID(ctx, uid)
		if err != nil {
			return err
		}
		for _, item := range items {
			resp.InsuranceIds = append(resp.InsuranceIds, item.ID.String())
		}
		return nil
	})

	g.Go(func() error {
		items, err := u.inv.GetInvestmentByUserID(ctx, uid)
		if err != nil {
			return err
		}
		for _, item := range items {
			resp.InvestmentIds = append(resp.InvestmentIds, item.ID.String())
		}
		return nil
	})

	g.Go(func() error {
		items, err := u.la.GetLandByUserID(ctx, uid)
		if err != nil {
			return err
		}
		for _, item := range items {
			resp.LandIds = append(resp.LandIds, item.ID.String())
		}
		return nil
	})

	g.Go(func() error {
		items, err := u.lia.GetLiabilityByUserID(ctx, uid)
		if err != nil {
			return err
		}
		for _, item := range items {
			resp.LiabilityIds = append(resp.LiabilityIds, item.ID.String())
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return resp, nil
}

func (u *AssetUsecase) GetAllAssets(ctx context.Context, req *pb.GetAllAssetsRequest) (*pb.GetAllAssetsResponse, error) {
	uid, _ := uuid.Parse(req.UserId)
	assets, err := u.r.GetAllAssets(ctx, uid)
	if err != nil {
		return nil, err
	}

	pbAssets := mapper.ToAssetSummaryProtoList(assets)
	return &pb.GetAllAssetsResponse{
		Assets: pbAssets,
	}, nil
}

func (u *AssetUsecase) GetNetWorth(ctx context.Context, req *pb.GetNetWorthRequest) (*pb.GetNetWorthResponse, error) {
	uid, _ := uuid.Parse(req.UserId)

	count, err := u.r.GetAssetCount(ctx, uid)
	if err != nil {
		return nil, err
	}

	overview, err := u.r.GetNetWorthOverview(ctx, uid)
	if err != nil {
		return nil, err
	}

	netWorth := overview.TotalAssets - overview.TotalLiabilities
	return &pb.GetNetWorthResponse{
		ItemCount:        count,
		AssetsValue:      overview.TotalAssets,
		LiabilitiesValue: overview.TotalLiabilities,
		NetWorth:         netWorth,
	}, nil
}
