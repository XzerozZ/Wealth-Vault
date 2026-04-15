package usecase_test

import (
	"context"
	"errors"
	"testing"
	"wealth-vault/asset-service/internal/domain"
	"wealth-vault/asset-service/internal/usecase"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	mock_repo "wealth-vault/asset-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAssetUsecase(t *testing.T) {
	repo := new(mock_repo.MockAssetRepository)
	uc := usecase.NewAssetUsecase(repo)

	ctx := context.Background()
	userID := uuid.New()
	assetID := uuid.New()

	t.Run("CheckExists_Success", func(t *testing.T) {
		req := &pb.CheckAssetRequest{
			Id:     assetID.String(),
			UserId: userID.String(),
			Type:   "account",
		}
		expectedName := "My Savings Account"

		repo.On("CheckExists", ctx, "account", assetID, userID).
			Return(expectedName, true, nil).Once()
		res, err := uc.CheckExists(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.True(t, res.Exists)
		assert.Equal(t, expectedName, res.Name)

		repo.AssertExpectations(t)
	})

	t.Run("GetAllAssetIDs_Success", func(t *testing.T) {
		req := &pb.GetMyAssetsRequest{UserId: userID.String()}
		mockMap := map[string][]string{
			"account":  {uuid.New().String(), uuid.New().String()},
			"building": {uuid.New().String()},
		}
		repo.On("GetAllAssetIDs", ctx, userID).Return(mockMap, nil).Once()

		res, err := uc.GetAllAssetIDs(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, mockMap["account"], res.AccountIds)
		assert.Equal(t, mockMap["building"], res.BuildingIds)
		assert.Empty(t, res.CashIds)
		repo.AssertExpectations(t)
	})

	t.Run("GetAllAssetIDs_InvalidUserID", func(t *testing.T) {
		req := &pb.GetMyAssetsRequest{UserId: "invalid-uuid"}
		res, err := uc.GetAllAssetIDs(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("GetAllAssets_Success", func(t *testing.T) {
		req := &pb.GetAllAssetsRequest{UserId: userID.String()}
		mockAssets := []domain.AssetSummary{
			{ID: assetID, Name: "My Home", Type: "building"},
		}

		mockLiabilities := []domain.AssetSummary{
			{ID: uuid.New(), Name: "Car Loan", Type: "liability"},
		}

		repo.On("GetAllAssets", ctx, userID).Return(mockAssets, mockLiabilities, nil).Once()
		res, err := uc.GetAllAssets(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res.Assets)
		repo.AssertExpectations(t)
	})

	t.Run("GetNetWorth_Success", func(t *testing.T) {
		req := &pb.GetNetWorthRequest{UserId: userID.String()}

		repo.On("GetAssetCount", ctx, userID).Return(int64(3), nil).Once()
		repo.On("GetNetWorthOverview", ctx, userID).Return(&domain.NetWorthOverview{
			TotalAssets:      500000,
			TotalLiabilities: 100000,
		}, nil).Once()

		res, err := uc.GetNetWorth(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), res.ItemCount)
		assert.Equal(t, float64(400000), res.NetWorth)
		repo.AssertExpectations(t)
	})

	t.Run("GetNetWorth_RepoError", func(t *testing.T) {
		req := &pb.GetNetWorthRequest{UserId: userID.String()}
		repo.On("GetAssetCount", ctx, userID).Return(int64(0), errors.New("db down")).Once()

		res, err := uc.GetNetWorth(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
