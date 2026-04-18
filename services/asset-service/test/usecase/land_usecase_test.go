package usecase_test

import (
	"context"
	"testing"
	"wealth-vault/asset-service/internal/domain"
	"wealth-vault/asset-service/internal/usecase"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	mock_client "wealth-vault/asset-service/test/mock/client"
	mock_helper "wealth-vault/asset-service/test/mock/helper"
	mock_repo "wealth-vault/asset-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLandUsecase(t *testing.T) {
	repo := new(mock_repo.MockLandRepository)
	assetHelper := new(mock_helper.MockAssetHelper)
	userClient := new(mock_client.MockUserClient)
	uc := usecase.NewLandUsecase(repo, assetHelper, userClient)

	ctx := context.Background()
	userID := uuid.New()
	landID := uuid.New()

	t.Run("CreateLand_Success", func(t *testing.T) {
		req := &pb.CreateLandRequest{
			UserId: userID.String(),
			Name:   "Green Valley Plot",
		}
		repo.On("CreateLand", ctx, mock.AnythingOfType("*domain.Land")).Return(nil).Once()

		res, err := uc.CreateLand(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, res.Land)
		repo.AssertExpectations(t)
	})

	t.Run("GetLand_Success", func(t *testing.T) {
		req := &pb.GetAssetRequest{UserId: userID.String()}
		lands := []*domain.Land{{ID: landID, Name: "Green Valley Plot"}}

		repo.On("GetLand", ctx, userID).Return(lands, nil).Once()

		res, err := uc.GetLand(ctx, req)
		assert.NoError(t, err)
		assert.Len(t, res.Land, 1)
		assert.True(t, res.Success)
	})

	t.Run("GetLandByID_Success", func(t *testing.T) {
		req := &pb.GetAssetByIDRequest{Id: landID.String()}
		repo.On("GetLandByID", ctx, landID).Return(&domain.Land{ID: landID}, nil).Once()

		res, err := uc.GetLandByID(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assert.Equal(t, landID.String(), res.Land.Id)
	})

	t.Run("UpdateLand_Success", func(t *testing.T) {
		buildingID := uuid.New()
		req := &pb.UpdateLandRequest{
			Id: landID.String(),
			Land: &pb.Land{
				UserId: userID.String(),
				Name:   "Updated Land Name",
			},
			BuildingIds: []string{buildingID.String()},
		}
		existingLand := &domain.Land{ID: landID, UserID: userID}

		repo.On("GetLandByID", ctx, landID).Return(existingLand, nil).Once()

		assetHelper.On("SyncFiles", ctx, mock.MatchedBy(func(p domain.FileSyncParams) bool {
			return p.EntityID == landID && p.EntityType == "land"
		})).Return(nil).Once()

		repo.On("UpdateLand", ctx, mock.Anything, []uuid.UUID{buildingID}, mock.Anything).
			Return(existingLand, nil).Once()

		res, err := uc.UpdateLand(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assetHelper.AssertExpectations(t)
	})

	t.Run("DeleteLand_Success", func(t *testing.T) {
		req := &pb.DeleteAssetRequest{Id: landID.String(), UserId: userID.String()}

		repo.On("GetLandByID", ctx, landID).Return(&domain.Land{ID: landID}, nil).Once()
		repo.On("SoftDeleteLand", ctx, landID, userID).Return(nil).Once()

		res, err := uc.DeleteLand(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
	})

	t.Run("CleanupExpiredLands_Success", func(t *testing.T) {
		expiredLands := []domain.Land{
			{ID: landID, Files: []domain.FileAssociate{{ID: uuid.New()}}},
		}

		repo.On("GetExpiredLand", ctx, mock.AnythingOfType("time.Time")).Return(expiredLands, nil).Once()

		assetHelper.On("CleanupResource", ctx, landID, expiredLands[0].Files, mock.Anything).
			Run(func(args mock.Arguments) {
				callback := args.Get(3).(func(uuid.UUID) error)
				callback(landID)
			}).Return().Once()

		repo.On("HardDeleteLand", ctx, landID).Return(nil).Once()

		err := uc.CleanupExpiredLands(ctx)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
		assetHelper.AssertExpectations(t)
	})

	t.Run("UpdateLand_MissingData", func(t *testing.T) {
		req := &pb.UpdateLandRequest{Id: landID.String(), Land: nil}
		res, err := uc.UpdateLand(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, "land data is required", err.Error())
		assert.Nil(t, res)
	})
}
