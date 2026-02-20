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

func TestCashUsecase(t *testing.T) {
	repo := new(mock_repo.MockCashRepository)
	assetHelper := new(mock_helper.MockAssetHelper)
	uc := usecase.NewCashUsecase(repo, assetHelper)

	ctx := context.Background()
	userID := uuid.New()
	cashID := uuid.New()

	t.Run("CreateCash_Success", func(t *testing.T) {
		req := &pb.CreateCashRequest{
			UserId: userID.String(),
			Name:   "Emergency Fund",
			Amount: 50000,
		}
		repo.On("CreateCash", ctx, mock.AnythingOfType("*domain.Cash")).Return(nil).Once()

		res, err := uc.CreateCash(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, res.Cash)
		repo.AssertExpectations(t)
	})

	t.Run("GetCash_Success", func(t *testing.T) {
		req := &pb.GetAssetRequest{UserId: userID.String()}
		cashes := []*domain.Cash{{ID: cashID, Name: "Emergency Fund"}}

		repo.On("GetCash", ctx, userID).Return(cashes, nil).Once()

		res, err := uc.GetCash(ctx, req)
		assert.NoError(t, err)
		assert.Len(t, res.Cash, 1)
		assert.True(t, res.Success)
	})

	t.Run("GetCashByID_Success", func(t *testing.T) {
		req := &pb.GetAssetByIDRequest{Id: cashID.String()}
		repo.On("GetCashByID", ctx, cashID).Return(&domain.Cash{ID: cashID}, nil).Once()

		res, err := uc.GetCashByID(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assert.Equal(t, cashID.String(), res.Cash.Id)
	})

	t.Run("UpdateCash_Success", func(t *testing.T) {
		req := &pb.UpdateCashRequest{
			Id: cashID.String(),
			Cash: &pb.Cash{
				UserId: userID.String(),
				Name:   "Updated Emergency Fund",
			},
			DeleteFileIds: []string{uuid.New().String()},
		}
		existingCash := &domain.Cash{ID: cashID, UserID: userID}

		repo.On("GetCashByID", ctx, cashID).Return(existingCash, nil).Once()

		assetHelper.On("SyncFiles", ctx, mock.MatchedBy(func(p domain.FileSyncParams) bool {
			return p.EntityID == cashID && p.EntityType == "cash"
		})).Return(nil).Once()

		repo.On("UpdateCash", ctx, mock.Anything).Return(existingCash, nil).Once()

		res, err := uc.UpdateCash(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assetHelper.AssertExpectations(t)
	})

	t.Run("DeleteCash_Success", func(t *testing.T) {
		req := &pb.DeleteAssetRequest{Id: cashID.String(), UserId: userID.String()}

		repo.On("GetCashByID", ctx, cashID).Return(&domain.Cash{ID: cashID}, nil).Once()
		repo.On("SoftDeleteCash", ctx, cashID, userID).Return(nil).Once()

		res, err := uc.DeleteCash(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
	})

	t.Run("CleanupExpiredCashes_Success", func(t *testing.T) {
		expiredCashes := []domain.Cash{
			{ID: cashID, Files: []domain.FileAssociate{{ID: uuid.New()}}},
		}

		repo.On("GetExpiredCash", ctx, mock.AnythingOfType("time.Time")).Return(expiredCashes, nil).Once()

		assetHelper.On("CleanupResource", ctx, cashID, expiredCashes[0].Files, mock.Anything).
			Run(func(args mock.Arguments) {
				callback := args.Get(3).(func(uuid.UUID) error)
				callback(cashID)
			}).Return().Once()

		repo.On("HardDeleteCash", ctx, cashID).Return(nil).Once()

		err := uc.CleanupExpiredCashes(ctx)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
		assetHelper.AssertExpectations(t)
	})

	t.Run("UpdateCash_InvalidID", func(t *testing.T) {
		req := &pb.UpdateCashRequest{Id: "invalid-uuid", Cash: &pb.Cash{UserId: userID.String()}}
		res, err := uc.UpdateCash(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("UpdateCash_MissingData", func(t *testing.T) {
		req := &pb.UpdateCashRequest{Id: cashID.String(), Cash: nil}
		res, err := uc.UpdateCash(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, "cash data is required", err.Error())
		assert.Nil(t, res)
	})
}
