package client

import (
	"context"

	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockAssetClient struct {
	mock.Mock
}

// CreateAccount implements pb.AssetServiceClient.
func (m *MockAssetClient) CreateAccount(ctx context.Context, in *assetPb.CreateAccountRequest, opts ...grpc.CallOption) (*assetPb.AccountResponse, error) {
	panic("unimplemented")
}

// CreateBuilding implements pb.AssetServiceClient.
func (m *MockAssetClient) CreateBuilding(ctx context.Context, in *assetPb.CreateBuildingRequest, opts ...grpc.CallOption) (*assetPb.BuildingResponse, error) {
	panic("unimplemented")
}

// CreateCash implements pb.AssetServiceClient.
func (m *MockAssetClient) CreateCash(ctx context.Context, in *assetPb.CreateCashRequest, opts ...grpc.CallOption) (*assetPb.CashResponse, error) {
	panic("unimplemented")
}

// CreateInsurance implements pb.AssetServiceClient.
func (m *MockAssetClient) CreateInsurance(ctx context.Context, in *assetPb.CreateInsuranceRequest, opts ...grpc.CallOption) (*assetPb.InsuranceResponse, error) {
	panic("unimplemented")
}

// CreateInvestment implements pb.AssetServiceClient.
func (m *MockAssetClient) CreateInvestment(ctx context.Context, in *assetPb.CreateInvestmentRequest, opts ...grpc.CallOption) (*assetPb.InvestmentResponse, error) {
	panic("unimplemented")
}

// CreateLand implements pb.AssetServiceClient.
func (m *MockAssetClient) CreateLand(ctx context.Context, in *assetPb.CreateLandRequest, opts ...grpc.CallOption) (*assetPb.LandResponse, error) {
	panic("unimplemented")
}

// CreateLiability implements pb.AssetServiceClient.
func (m *MockAssetClient) CreateLiability(ctx context.Context, in *assetPb.CreateLiabilityRequest, opts ...grpc.CallOption) (*assetPb.LiabilityResponse, error) {
	panic("unimplemented")
}

// DeleteAccount implements pb.AssetServiceClient.
func (m *MockAssetClient) DeleteAccount(ctx context.Context, in *assetPb.DeleteAssetRequest, opts ...grpc.CallOption) (*assetPb.DeleteAssetResponse, error) {
	panic("unimplemented")
}

// DeleteBuilding implements pb.AssetServiceClient.
func (m *MockAssetClient) DeleteBuilding(ctx context.Context, in *assetPb.DeleteAssetRequest, opts ...grpc.CallOption) (*assetPb.DeleteAssetResponse, error) {
	panic("unimplemented")
}

// DeleteCash implements pb.AssetServiceClient.
func (m *MockAssetClient) DeleteCash(ctx context.Context, in *assetPb.DeleteAssetRequest, opts ...grpc.CallOption) (*assetPb.DeleteAssetResponse, error) {
	panic("unimplemented")
}

// DeleteInsurance implements pb.AssetServiceClient.
func (m *MockAssetClient) DeleteInsurance(ctx context.Context, in *assetPb.DeleteAssetRequest, opts ...grpc.CallOption) (*assetPb.DeleteAssetResponse, error) {
	panic("unimplemented")
}

// DeleteInvestment implements pb.AssetServiceClient.
func (m *MockAssetClient) DeleteInvestment(ctx context.Context, in *assetPb.DeleteAssetRequest, opts ...grpc.CallOption) (*assetPb.DeleteAssetResponse, error) {
	panic("unimplemented")
}

// DeleteLand implements pb.AssetServiceClient.
func (m *MockAssetClient) DeleteLand(ctx context.Context, in *assetPb.DeleteAssetRequest, opts ...grpc.CallOption) (*assetPb.DeleteAssetResponse, error) {
	panic("unimplemented")
}

// DeleteLiability implements pb.AssetServiceClient.
func (m *MockAssetClient) DeleteLiability(ctx context.Context, in *assetPb.DeleteLiabilityRequest, opts ...grpc.CallOption) (*assetPb.DeleteLiabilityResponse, error) {
	panic("unimplemented")
}

// GetAccount implements pb.AssetServiceClient.
func (m *MockAssetClient) GetAccount(ctx context.Context, in *assetPb.GetAssetRequest, opts ...grpc.CallOption) (*assetPb.AccountArrayResponse, error) {
	panic("unimplemented")
}

// GetAccountByID implements pb.AssetServiceClient.
func (m *MockAssetClient) GetAccountByID(ctx context.Context, in *assetPb.GetAssetByIDRequest, opts ...grpc.CallOption) (*assetPb.AccountResponse, error) {
	panic("unimplemented")
}

// GetAllAssets implements pb.AssetServiceClient.
func (m *MockAssetClient) GetAllAssets(ctx context.Context, in *assetPb.GetAllAssetsRequest, opts ...grpc.CallOption) (*assetPb.GetAllAssetsResponse, error) {
	panic("unimplemented")
}

// GetBuilding implements pb.AssetServiceClient.
func (m *MockAssetClient) GetBuilding(ctx context.Context, in *assetPb.GetAssetRequest, opts ...grpc.CallOption) (*assetPb.BuildingArrayResponse, error) {
	panic("unimplemented")
}

// GetBuildingByID implements pb.AssetServiceClient.
func (m *MockAssetClient) GetBuildingByID(ctx context.Context, in *assetPb.GetAssetByIDRequest, opts ...grpc.CallOption) (*assetPb.BuildingResponse, error) {
	panic("unimplemented")
}

// GetCash implements pb.AssetServiceClient.
func (m *MockAssetClient) GetCash(ctx context.Context, in *assetPb.GetAssetRequest, opts ...grpc.CallOption) (*assetPb.CashArrayResponse, error) {
	panic("unimplemented")
}

// GetCashByID implements pb.AssetServiceClient.
func (m *MockAssetClient) GetCashByID(ctx context.Context, in *assetPb.GetAssetByIDRequest, opts ...grpc.CallOption) (*assetPb.CashResponse, error) {
	panic("unimplemented")
}

// GetInsurance implements pb.AssetServiceClient.
func (m *MockAssetClient) GetInsurance(ctx context.Context, in *assetPb.GetAssetRequest, opts ...grpc.CallOption) (*assetPb.InsuranceArrayResponse, error) {
	panic("unimplemented")
}

// GetInsuranceByID implements pb.AssetServiceClient.
func (m *MockAssetClient) GetInsuranceByID(ctx context.Context, in *assetPb.GetAssetByIDRequest, opts ...grpc.CallOption) (*assetPb.InsuranceResponse, error) {
	panic("unimplemented")
}

// GetInvestment implements pb.AssetServiceClient.
func (m *MockAssetClient) GetInvestment(ctx context.Context, in *assetPb.GetAssetRequest, opts ...grpc.CallOption) (*assetPb.InvestmentArrayResponse, error) {
	panic("unimplemented")
}

// GetInvestmentByID implements pb.AssetServiceClient.
func (m *MockAssetClient) GetInvestmentByID(ctx context.Context, in *assetPb.GetAssetByIDRequest, opts ...grpc.CallOption) (*assetPb.InvestmentResponse, error) {
	panic("unimplemented")
}

// GetLand implements pb.AssetServiceClient.
func (m *MockAssetClient) GetLand(ctx context.Context, in *assetPb.GetAssetRequest, opts ...grpc.CallOption) (*assetPb.LandArrayResponse, error) {
	panic("unimplemented")
}

// GetLandByID implements pb.AssetServiceClient.
func (m *MockAssetClient) GetLandByID(ctx context.Context, in *assetPb.GetAssetByIDRequest, opts ...grpc.CallOption) (*assetPb.LandResponse, error) {
	panic("unimplemented")
}

// GetLiability implements pb.AssetServiceClient.
func (m *MockAssetClient) GetLiability(ctx context.Context, in *assetPb.GetLiabilityRequest, opts ...grpc.CallOption) (*assetPb.LiabilityArrayResponse, error) {
	panic("unimplemented")
}

// GetLiabilityByID implements pb.AssetServiceClient.
func (m *MockAssetClient) GetLiabilityByID(ctx context.Context, in *assetPb.GetLiabilityByIDRequest, opts ...grpc.CallOption) (*assetPb.LiabilityResponse, error) {
	panic("unimplemented")
}

// GetNetWorth implements pb.AssetServiceClient.
func (m *MockAssetClient) GetNetWorth(ctx context.Context, in *assetPb.GetNetWorthRequest, opts ...grpc.CallOption) (*assetPb.GetNetWorthResponse, error) {
	panic("unimplemented")
}

// UpdateAccount implements pb.AssetServiceClient.
func (m *MockAssetClient) UpdateAccount(ctx context.Context, in *assetPb.UpdateAccountRequest, opts ...grpc.CallOption) (*assetPb.AccountResponse, error) {
	panic("unimplemented")
}

// UpdateBuilding implements pb.AssetServiceClient.
func (m *MockAssetClient) UpdateBuilding(ctx context.Context, in *assetPb.UpdateBuildingRequest, opts ...grpc.CallOption) (*assetPb.BuildingResponse, error) {
	panic("unimplemented")
}

// UpdateCash implements pb.AssetServiceClient.
func (m *MockAssetClient) UpdateCash(ctx context.Context, in *assetPb.UpdateCashRequest, opts ...grpc.CallOption) (*assetPb.CashResponse, error) {
	panic("unimplemented")
}

// UpdateInsurance implements pb.AssetServiceClient.
func (m *MockAssetClient) UpdateInsurance(ctx context.Context, in *assetPb.UpdateInsuranceRequest, opts ...grpc.CallOption) (*assetPb.InsuranceResponse, error) {
	panic("unimplemented")
}

// UpdateInvestment implements pb.AssetServiceClient.
func (m *MockAssetClient) UpdateInvestment(ctx context.Context, in *assetPb.UpdateInvestmentRequest, opts ...grpc.CallOption) (*assetPb.InvestmentResponse, error) {
	panic("unimplemented")
}

// UpdateLand implements pb.AssetServiceClient.
func (m *MockAssetClient) UpdateLand(ctx context.Context, in *assetPb.UpdateLandRequest, opts ...grpc.CallOption) (*assetPb.LandResponse, error) {
	panic("unimplemented")
}

// UpdateLiability implements pb.AssetServiceClient.
func (m *MockAssetClient) UpdateLiability(ctx context.Context, in *assetPb.UpdateLiabilityRequest, opts ...grpc.CallOption) (*assetPb.LiabilityResponse, error) {
	panic("unimplemented")
}

func (m *MockAssetClient) CheckAssetExists(
	ctx context.Context,
	in *assetPb.CheckAssetRequest,
	opts ...grpc.CallOption,
) (*assetPb.CheckAssetResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*assetPb.CheckAssetResponse), args.Error(1)
}

func (m *MockAssetClient) GetAllAssetIDs(ctx context.Context, req *assetPb.GetMyAssetsRequest, opts ...grpc.CallOption) (*assetPb.GetMyAssetsResponse, error) {
	args := m.Called(ctx, req, opts)
	return args.Get(0).(*assetPb.GetMyAssetsResponse), args.Error(1)
}

func (m *MockAssetClient) GetBatchBuilding(ctx context.Context, req *assetPb.GetBatchIdsRequest, opts ...grpc.CallOption) (*assetPb.BuildingArrayResponse, error) {
	args := m.Called(ctx, req, opts)
	return args.Get(0).(*assetPb.BuildingArrayResponse), args.Error(1)
}

func (m *MockAssetClient) GetBatchLand(ctx context.Context, req *assetPb.GetBatchIdsRequest, opts ...grpc.CallOption) (*assetPb.LandArrayResponse, error) {
	args := m.Called(ctx, req, opts)
	return args.Get(0).(*assetPb.LandArrayResponse), args.Error(1)
}

func (m *MockAssetClient) GetBatchAccount(ctx context.Context, req *assetPb.GetBatchIdsRequest, opts ...grpc.CallOption) (*assetPb.AccountArrayResponse, error) {
	args := m.Called(ctx, req, opts)
	return args.Get(0).(*assetPb.AccountArrayResponse), args.Error(1)
}

func (m *MockAssetClient) GetBatchCash(ctx context.Context, req *assetPb.GetBatchIdsRequest, opts ...grpc.CallOption) (*assetPb.CashArrayResponse, error) {
	args := m.Called(ctx, req, opts)
	return args.Get(0).(*assetPb.CashArrayResponse), args.Error(1)
}

func (m *MockAssetClient) GetBatchInsurance(ctx context.Context, req *assetPb.GetBatchIdsRequest, opts ...grpc.CallOption) (*assetPb.InsuranceArrayResponse, error) {
	args := m.Called(ctx, req, opts)
	return args.Get(0).(*assetPb.InsuranceArrayResponse), args.Error(1)
}

func (m *MockAssetClient) GetBatchInvestment(ctx context.Context, req *assetPb.GetBatchIdsRequest, opts ...grpc.CallOption) (*assetPb.InvestmentArrayResponse, error) {
	args := m.Called(ctx, req, opts)
	return args.Get(0).(*assetPb.InvestmentArrayResponse), args.Error(1)
}

func (m *MockAssetClient) GetBatchLiability(ctx context.Context, req *assetPb.GetBatchIdsRequest, opts ...grpc.CallOption) (*assetPb.LiabilityArrayResponse, error) {
	args := m.Called(ctx, req, opts)
	return args.Get(0).(*assetPb.LiabilityArrayResponse), args.Error(1)
}
