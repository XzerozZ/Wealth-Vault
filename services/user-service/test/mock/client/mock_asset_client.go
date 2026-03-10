package client

import (
	"context"

	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockAssetClient struct {
	assetPb.AssetServiceClient
	mock.Mock
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
