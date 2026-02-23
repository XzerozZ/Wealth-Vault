package usecase_test

import (
	"context"
	"errors"
	"testing"
	"wealth-vault/user-service/internal/domain"
	usecase "wealth-vault/user-service/internal/usecase/shareItem"
	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"
	mock_client "wealth-vault/user-service/test/mock/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShareItemUsecase_FetchAssetPreviews(t *testing.T) {
	t.Run("success multiple types", func(t *testing.T) {
		mockAsset := new(mock_client.MockAssetClient)

		uc := usecase.NewShareItemUsecase(nil, nil, nil, nil, mockAsset, nil, nil)

		items := []domain.SharedItemSummary{
			{EntityID: "b1", EntityType: usecase.AssetTypeBuilding},
			{EntityID: "a1", EntityType: usecase.AssetTypeAccount},
		}

		mockAsset.On("GetBatchBuilding", mock.Anything, mock.Anything, mock.Anything).
			Return(&assetPb.BuildingArrayResponse{
				Building: []*assetPb.Building{
					{Id: "b1"},
				},
			}, nil)

		mockAsset.On("GetBatchAccount", mock.Anything, mock.Anything, mock.Anything).
			Return(&assetPb.AccountArrayResponse{
				Account: []*assetPb.Account{
					{Id: "a1"},
				},
			}, nil)

		result, err := uc.FetchAssetPreviews(context.Background(), items)

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		mockAsset.AssertExpectations(t)
	})

	t.Run("error in one type should not fail all", func(t *testing.T) {
		mockAsset := new(mock_client.MockAssetClient)

		uc := usecase.NewShareItemUsecase(nil, nil, nil, nil, mockAsset, nil, nil)

		items := []domain.SharedItemSummary{
			{EntityID: "b1", EntityType: usecase.AssetTypeBuilding},
		}

		mockAsset.On("GetBatchBuilding", mock.Anything, mock.Anything, mock.Anything).
			Return((*assetPb.BuildingArrayResponse)(nil), errors.New("grpc error"))

		result, err := uc.FetchAssetPreviews(context.Background(), items)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("empty input", func(t *testing.T) {
		mockAsset := new(mock_client.MockAssetClient)

		uc := usecase.NewShareItemUsecase(nil, nil, nil, nil, mockAsset, nil, nil)

		result, err := uc.FetchAssetPreviews(context.Background(), []domain.SharedItemSummary{})

		assert.NoError(t, err)
		assert.Empty(t, result)

		mockAsset.AssertExpectations(t)
	})

	t.Run("lowercase type handling", func(t *testing.T) {
		mockAsset := new(mock_client.MockAssetClient)
		uc := usecase.NewShareItemUsecase(nil, nil, nil, nil, mockAsset, nil, nil)

		items := []domain.SharedItemSummary{
			{EntityID: "b1", EntityType: "BUILDING"},
		}

		mockAsset.On("GetBatchBuilding", mock.Anything, mock.Anything, mock.Anything).
			Return(&assetPb.BuildingArrayResponse{
				Building: []*assetPb.Building{
					{
						Id:   "b1",
						Name: "ตึกทดสอบ",
						Location: &assetPb.Location{
							District: "เขต",
							Province: "จังหวัด",
						},
					},
				},
			}, nil)

		result, err := uc.FetchAssetPreviews(context.Background(), items)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		mockAsset.AssertExpectations(t)
	})
}
