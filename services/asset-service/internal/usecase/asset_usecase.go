package usecase

import (
	"context"
	repo "wealth-vault/asset-service/internal/repository/interface"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils"
	"wealth-vault/asset-service/pkg/utils/mapper"
)

type AssetUsecase struct {
	r repo.AssetRepository
}

func NewAssetUsecase(r repo.AssetRepository) *AssetUsecase {
	return &AssetUsecase{r: r}
}

func (u *AssetUsecase) CheckExists(ctx context.Context, req *pb.CheckAssetRequest) (*pb.CheckAssetResponse, error) {
	id, uid, err := utils.ValidateIDs(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}

	name, exist, err := u.r.CheckExists(ctx, req.Type, id, uid)
	if err != nil {
		return nil, err
	}

	return &pb.CheckAssetResponse{
		Exists: exist,
		Name:   name,
	}, nil
}

func (u *AssetUsecase) GetAllAssetIDs(ctx context.Context, req *pb.GetMyAssetsRequest) (*pb.GetMyAssetsResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	assetMap, err := u.r.GetAllAssetIDs(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &pb.GetMyAssetsResponse{
		AccountIds:    assetMap["account"],
		BuildingIds:   assetMap["building"],
		CashIds:       assetMap["cash"],
		InsuranceIds:  assetMap["insurance"],
		InvestmentIds: assetMap["investment"],
		LandIds:       assetMap["land"],
		LiabilityIds:  assetMap["liability"],
	}, nil
}

func (u *AssetUsecase) GetAllAssetsSelection(ctx context.Context, req *pb.GetAllAssetsRequest) (*pb.GetAllAssetsResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	assets, lias, err := u.r.GetAllAssetSelection(ctx, uid)
	if err != nil {
		return nil, err
	}

	pbAssets := mapper.ToAssetSummaryProtoList(assets)
	pbLias := mapper.ToAssetSummaryProtoList(lias)
	return &pb.GetAllAssetsResponse{
		Assets:      pbAssets,
		Liabilities: pbLias,
	}, nil
}

func (u *AssetUsecase) GetAllAssets(ctx context.Context, req *pb.GetAllAssetsRequest) (*pb.GetAllAssetsResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

	assets, lias, err := u.r.GetAllAssets(ctx, uid)
	if err != nil {
		return nil, err
	}

	pbAssets := mapper.ToAssetSummaryProtoList(assets)
	pbLias := mapper.ToAssetSummaryProtoList(lias)
	return &pb.GetAllAssetsResponse{
		Assets:      pbAssets,
		Liabilities: pbLias,
	}, nil
}

func (u *AssetUsecase) GetNetWorth(ctx context.Context, req *pb.GetNetWorthRequest) (*pb.GetNetWorthResponse, error) {
	uid, err := utils.ParseID(req.UserId)
	if err != nil {
		return nil, err
	}

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
