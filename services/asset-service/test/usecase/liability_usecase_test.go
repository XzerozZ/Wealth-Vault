package usecase_test

import (
	"context"
	"testing"
	"wealth-vault/asset-service/internal/domain"
	"wealth-vault/asset-service/internal/usecase"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	mock_helper "wealth-vault/asset-service/test/mock/helper"
	mock_repo "wealth-vault/asset-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLiabilityUsecase(t *testing.T) {
	repo := new(mock_repo.MockLiabilityRepository)
	assetHelper := new(mock_helper.MockAssetHelper)
	uc := usecase.NewLiabilityUsecase(repo, assetHelper)

	ctx := context.Background()
	userID := uuid.New()
	liaID := uuid.New()

	t.Run("CreateLiability_Success", func(t *testing.T) {
		req := &pb.CreateLiabilityRequest{
			UserId: userID.String(),
			Name:   "Car Loan",
		}
		repo.On("CreateLiability", ctx, mock.AnythingOfType("*domain.Liability")).Return(nil).Once()

		res, err := uc.CreateLiability(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assert.NotNil(t, res.Liability)
		repo.AssertExpectations(t)
	})

	t.Run("GetLiability_Success", func(t *testing.T) {
		req := &pb.GetLiabilityRequest{UserId: userID.String()}
		lias := []*domain.Liability{{ID: liaID, Name: "Car Loan"}}

		repo.On("GetLiability", ctx, userID).Return(lias, nil).Once()

		res, err := uc.GetLiability(ctx, req)
		assert.NoError(t, err)
		assert.Len(t, res.Liability, 1)
		assert.True(t, res.Success)
	})

	t.Run("GetLiabilityByID_Success", func(t *testing.T) {
		req := &pb.GetLiabilityByIDRequest{Id: liaID.String()}
		repo.On("GetLiabilityByID", ctx, liaID).Return(&domain.Liability{ID: liaID}, nil).Once()

		res, err := uc.GetLiabilityByID(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assert.Equal(t, liaID.String(), res.Liability.Id)
	})

	t.Run("UpdateLiability_Success", func(t *testing.T) {
		req := &pb.UpdateLiabilityRequest{
			Id: liaID.String(),
			Liability: &pb.Liability{
				CreatedBy: userID.String(),
				Name:      "Updated Loan Name",
			},
		}
		existingLia := &domain.Liability{ID: liaID, UserID: userID}

		repo.On("GetLiabilityByID", ctx, liaID).Return(existingLia, nil).Once()

		assetHelper.On("SyncFiles", ctx, mock.MatchedBy(func(p domain.FileSyncParams) bool {
			return p.EntityID == liaID && p.EntityType == "liability"
		})).Return(nil).Once()

		repo.On("UpdateLiability", ctx, mock.Anything).Return(existingLia, nil).Once()

		res, err := uc.UpdateLiability(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assetHelper.AssertExpectations(t)
	})

	t.Run("DeleteLiability_Success", func(t *testing.T) {
		req := &pb.DeleteLiabilityRequest{Id: liaID.String(), UserId: userID.String()}

		repo.On("GetLiabilityByID", ctx, liaID).Return(&domain.Liability{ID: liaID}, nil).Once()
		repo.On("SoftDeleteLiability", ctx, liaID, userID).Return(nil).Once()

		res, err := uc.DeleteLiability(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
	})

	t.Run("CleanupExpiredLiabilities_Success", func(t *testing.T) {
		expiredLias := []domain.Liability{
			{ID: liaID, Files: []domain.FileAssociate{{ID: uuid.New()}}},
		}

		repo.On("GetExpiredLiability", ctx, mock.AnythingOfType("time.Time")).Return(expiredLias, nil).Once()

		assetHelper.On("CleanupResource", ctx, liaID, expiredLias[0].Files, mock.Anything).
			Run(func(args mock.Arguments) {
				cb := args.Get(3).(func(uuid.UUID) error)
				cb(liaID)
			}).Return().Once()

		repo.On("HardDeleteLiability", ctx, liaID).Return(nil).Once()

		err := uc.CleanupExpiredLiabilities(ctx)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("UpdateLiability_MissingData", func(t *testing.T) {
		req := &pb.UpdateLiabilityRequest{Id: liaID.String(), Liability: nil}
		res, err := uc.UpdateLiability(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, "liability data is required", err.Error())
		assert.Nil(t, res)
	})
}
