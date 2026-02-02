package grpc

import (
	"context"
	"wealth-vault/asset-service/internal/usecase"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

type AssetGRPCHandler struct {
	pb.UnimplementedAssetServiceServer
	assetusecase usecase.AssetUsecase
	accusecase   usecase.AccountUsecase
	cashusecase  usecase.CashUsecase
	inusecase    usecase.InvestmentUsecase
	buildusecase usecase.BuildingUsecase
	landusecase  usecase.LandUsecase
	insusecase   usecase.InsuranceUsecase
	liausecase   usecase.LiabilityUsecase
}

func NewAssetGRPCHandler(
	as usecase.AssetUsecase,
	ac usecase.AccountUsecase,
	cu usecase.CashUsecase,
	iu usecase.InvestmentUsecase,
	bu usecase.BuildingUsecase,
	la usecase.LandUsecase,
	in usecase.InsuranceUsecase,
	lu usecase.LiabilityUsecase,
) *AssetGRPCHandler {
	return &AssetGRPCHandler{
		assetusecase: as,
		accusecase:   ac,
		cashusecase:  cu,
		inusecase:    iu,
		buildusecase: bu,
		landusecase:  la,
		insusecase:   in,
		liausecase:   lu,
	}
}

func (h *AssetGRPCHandler) CheckAssetExists(ctx context.Context, req *pb.CheckAssetRequest) (*pb.CheckAssetResponse, error) {
	res, err := h.assetusecase.CheckExists(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.AccountResponse, error) {
	res, err := h.accusecase.CreateAccount(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) CreateCash(ctx context.Context, req *pb.CreateCashRequest) (*pb.CashResponse, error) {
	res, err := h.cashusecase.CreateCash(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) CreateInvestment(ctx context.Context, req *pb.CreateInvestmentRequest) (*pb.InvestmentResponse, error) {
	res, err := h.inusecase.CreateInvestment(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) CreateLiability(ctx context.Context, req *pb.CreateLiabilityRequest) (*pb.LiabilityResponse, error) {
	res, err := h.liausecase.CreateLiability(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (h *AssetGRPCHandler) CreateBuilding(ctx context.Context, req *pb.CreateBuildingRequest) (*pb.BuildingResponse, error) {
	res, err := h.buildusecase.CreateBuilding(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) CreateLand(ctx context.Context, req *pb.CreateLandRequest) (*pb.LandResponse, error) {
	res, err := h.landusecase.CreateLand(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) CreateInsurance(ctx context.Context, req *pb.CreateInsuranceRequest) (*pb.InsuranceResponse, error) {
	res, err := h.insusecase.CreateInsurance(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetAccount(ctx context.Context, req *pb.GetAssetRequest) (*pb.AccountArrayResponse, error) {
	res, err := h.accusecase.GetAccount(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetAccountByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.AccountResponse, error) {
	res, err := h.accusecase.GetAccountByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetCash(ctx context.Context, req *pb.GetAssetRequest) (*pb.CashArrayResponse, error) {
	res, err := h.cashusecase.GetCash(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetCashByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.CashResponse, error) {
	res, err := h.cashusecase.GetCashByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetInvestment(ctx context.Context, req *pb.GetAssetRequest) (*pb.InvestmentArrayResponse, error) {
	res, err := h.inusecase.GetInvestment(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetInvestmentByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.InvestmentResponse, error) {
	res, err := h.inusecase.GetInvestmentByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetBuilding(ctx context.Context, req *pb.GetAssetRequest) (*pb.BuildingArrayResponse, error) {
	res, err := h.buildusecase.GetBuilding(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetBuildingByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.BuildingResponse, error) {
	res, err := h.buildusecase.GetBuildingByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetLand(ctx context.Context, req *pb.GetAssetRequest) (*pb.LandArrayResponse, error) {
	res, err := h.landusecase.GetLand(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetLandByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.LandResponse, error) {
	res, err := h.landusecase.GetLandByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetInsurance(ctx context.Context, req *pb.GetAssetRequest) (*pb.InsuranceArrayResponse, error) {
	res, err := h.insusecase.GetInsurance(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetInsuranceByID(ctx context.Context, req *pb.GetAssetByIDRequest) (*pb.InsuranceResponse, error) {
	res, err := h.insusecase.GetInsuranceByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetLiability(ctx context.Context, req *pb.GetLiabilityRequest) (*pb.LiabilityArrayResponse, error) {
	lias, err := h.liausecase.GetLiability(ctx, req)
	if err != nil {
		return nil, err
	}

	return lias, nil
}

func (h *AssetGRPCHandler) GetLiabilityByID(ctx context.Context, req *pb.GetLiabilityByIDRequest) (*pb.LiabilityResponse, error) {
	lia, err := h.liausecase.GetLiabilityByID(ctx, req)
	if err != nil {
		return nil, err
	}

	return lia, nil
}

func (h *AssetGRPCHandler) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.AccountResponse, error) {
	res, err := h.accusecase.UpdateAccount(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) UpdateCash(ctx context.Context, req *pb.UpdateCashRequest) (*pb.CashResponse, error) {
	res, err := h.cashusecase.UpdateCash(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) UpdateInvestment(ctx context.Context, req *pb.UpdateInvestmentRequest) (*pb.InvestmentResponse, error) {
	res, err := h.inusecase.UpdateInvestment(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) UpdateBuilding(ctx context.Context, req *pb.UpdateBuildingRequest) (*pb.BuildingResponse, error) {
	res, err := h.buildusecase.UpdateBuilding(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) UpdateLand(ctx context.Context, req *pb.UpdateLandRequest) (*pb.LandResponse, error) {
	res, err := h.landusecase.UpdateLand(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) UpdateInsurance(ctx context.Context, req *pb.UpdateInsuranceRequest) (*pb.InsuranceResponse, error) {
	res, err := h.insusecase.UpdateInsurance(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) UpdateLiability(ctx context.Context, req *pb.UpdateLiabilityRequest) (*pb.LiabilityResponse, error) {
	lia, err := h.liausecase.UpdateLiability(ctx, req)
	if err != nil {
		return nil, err
	}

	return lia, nil
}

func (h *AssetGRPCHandler) DeleteAccount(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	res, err := h.accusecase.DeleteAccount(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) DeleteCash(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	res, err := h.cashusecase.DeleteCash(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (h *AssetGRPCHandler) DeleteInvestment(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	res, err := h.inusecase.DeleteInvestment(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) DeleteBuilding(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	res, err := h.buildusecase.DeleteBuilding(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) DeleteLand(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	res, err := h.landusecase.DeleteLand(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) DeleteInsurance(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	res, err := h.insusecase.DeleteInsurance(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) DeleteLiability(ctx context.Context, req *pb.DeleteLiabilityRequest) (*pb.DeleteLiabilityResponse, error) {
	lia, err := h.liausecase.DeleteLiability(ctx, req)
	if err != nil {
		return nil, err
	}

	return lia, nil
}

func (h *AssetGRPCHandler) GetBatchAccount(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.AccountArrayResponse, error) {
	res, err := h.accusecase.GetAccountByIDs(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetBatchBuilding(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.BuildingArrayResponse, error) {
	res, err := h.buildusecase.GetBuildingByIDs(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetBatchCash(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.CashArrayResponse, error) {
	res, err := h.cashusecase.GetCashByIDs(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetBatchInsurance(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.InsuranceArrayResponse, error) {
	res, err := h.insusecase.GetInsuranceByIDs(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetBatchInvestment(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.InvestmentArrayResponse, error) {
	res, err := h.inusecase.GetInvestmentByIDs(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetBatchLand(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.LandArrayResponse, error) {
	res, err := h.landusecase.GetLandByIDs(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *AssetGRPCHandler) GetBatchLiability(ctx context.Context, req *pb.GetBatchIdsRequest) (*pb.LiabilityArrayResponse, error) {
	res, err := h.liausecase.GetLiabilityByIDs(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
