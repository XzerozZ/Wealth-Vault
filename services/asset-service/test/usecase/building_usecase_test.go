package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"wealth-vault/asset-service/internal/domain"
	"wealth-vault/asset-service/internal/usecase"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	mock_helper "wealth-vault/asset-service/test/mock/helper"
	mock_repo "wealth-vault/asset-service/test/mock/repository"
)

func TestBuildingUsecase(t *testing.T) {
	repo := new(mock_repo.MockBuildingRepository)
	assetHelper := new(mock_helper.MockAssetHelper)
	uc := usecase.NewBuildingUsecase(repo, assetHelper)

	ctx := context.Background()
	userID := uuid.New()
	buildingID := uuid.New()

	t.Run("CreateBuilding_Success", func(t *testing.T) {
		req := &pb.CreateBuildingRequest{
			UserId:   userID.String(),
			Name:     "Modern Condo",
			Location: &pb.Location{Province: "Bangkok"},
		}
		repo.On("CreateBuilding", ctx, mock.AnythingOfType("*domain.Building")).Return(nil).Once()

		res, err := uc.CreateBuilding(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, res.Building)
		repo.AssertExpectations(t)
	})

	t.Run("GetBuilding_Success", func(t *testing.T) {
		req := &pb.GetAssetRequest{UserId: userID.String()}
		buildings := []*domain.Building{{ID: buildingID, Name: "Condo A"}}

		repo.On("GetBuilding", ctx, userID).Return(buildings, nil).Once()

		res, err := uc.GetBuilding(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assert.Len(t, res.Building, 1)
	})

	t.Run("GetBuildingByID_Success", func(t *testing.T) {
		req := &pb.GetAssetByIDRequest{Id: buildingID.String()}
		repo.On("GetBuildingByID", ctx, buildingID).Return(&domain.Building{ID: buildingID}, nil).Once()

		res, err := uc.GetBuildingByID(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, buildingID.String(), res.Building.Id)
	})

	t.Run("UpdateBuilding_Success", func(t *testing.T) {
		req := &pb.UpdateBuildingRequest{
			Id:       buildingID.String(),
			Building: &pb.Building{UserId: userID.String(), Name: "Updated Condo"},
		}
		existing := &domain.Building{ID: buildingID, UserID: userID}

		repo.On("GetBuildingByID", ctx, buildingID).Return(existing, nil).Once()
		assetHelper.On("SyncFiles", ctx, mock.MatchedBy(func(p domain.FileSyncParams) bool {
			return p.EntityID == buildingID && p.EntityType == "building"
		})).Return(nil).Once()

		repo.On("UpdateBuilding", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(existing, nil).Once()

		res, err := uc.UpdateBuilding(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assetHelper.AssertExpectations(t)
	})

	t.Run("UpdateBuilding_NoDataError", func(t *testing.T) {
		res, err := uc.UpdateBuilding(ctx, &pb.UpdateBuildingRequest{Building: nil})
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "required")
	})

	t.Run("DeleteBuilding_Success", func(t *testing.T) {
		req := &pb.DeleteAssetRequest{Id: buildingID.String(), UserId: userID.String()}

		repo.On("GetBuildingByID", ctx, buildingID).Return(&domain.Building{}, nil).Once()
		repo.On("SoftDeleteBuilding", ctx, buildingID, userID).Return(nil).Once()

		res, err := uc.DeleteBuilding(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
	})

	t.Run("CleanupExpiredBuildings_Success", func(t *testing.T) {
		expired := []domain.Building{{ID: buildingID, Files: []domain.FileAssociate{{ID: uuid.New()}}}}

		repo.On("GetExpiredBuilding", ctx, mock.AnythingOfType("time.Time")).Return(expired, nil).Once()

		assetHelper.On("CleanupResource", ctx, buildingID, expired[0].Files, mock.Anything).
			Run(func(args mock.Arguments) {
				callback := args.Get(3).(func(uuid.UUID) error)
				callback(buildingID)
			}).Return().Once()

		repo.On("HardDeleteBuilding", ctx, buildingID).Return(nil).Once()

		err := uc.CleanupExpiredBuildings(ctx)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
		assetHelper.AssertExpectations(t)
	})

	t.Run("CleanupExpiredBuildings_Empty", func(t *testing.T) {
		repo.On("GetExpiredBuilding", ctx, mock.Anything).Return([]domain.Building{}, nil).Once()

		err := uc.CleanupExpiredBuildings(ctx)
		assert.NoError(t, err)
	})
}
