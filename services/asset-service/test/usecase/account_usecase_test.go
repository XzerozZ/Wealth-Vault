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

func TestAccountUsecase(t *testing.T) {
	repo := new(mock_repo.MockAccountRepository)
	assetHelper := new(mock_helper.MockAssetHelper)
	uc := usecase.NewAccountUsecase(repo, assetHelper)

	ctx := context.Background()
	userID := uuid.New()
	accountID := uuid.New()

	t.Run("CreateAccount_Success", func(t *testing.T) {
		req := &pb.CreateAccountRequest{UserId: userID.String(), Name: "Savings"}
		repo.On("CreateAccount", ctx, mock.AnythingOfType("*domain.Account")).Return(nil).Once()

		res, err := uc.CreateAccount(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, res.Account)
		repo.AssertExpectations(t)
	})

	t.Run("GetAccount_Success", func(t *testing.T) {
		req := &pb.GetAssetRequest{
			UserId: userID.String(),
		}
		accounts := []*domain.Account{{ID: accountID, Name: "Savings"}}

		repo.On("GetAccount", ctx, userID).Return(accounts, nil).Once()

		res, err := uc.GetAccount(ctx, req)
		assert.NoError(t, err)
		assert.Len(t, res.Account, 1)
		assert.True(t, res.Success)
	})

	t.Run("GetAccountByID_Success", func(t *testing.T) {
		req := &pb.GetAssetByIDRequest{Id: accountID.String()}
		repo.On("GetAccountByID", ctx, accountID).Return(&domain.Account{ID: accountID}, nil).Once()

		res, err := uc.GetAccountByID(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
	})

	t.Run("UpdateAccount_Success", func(t *testing.T) {
		req := &pb.UpdateAccountRequest{
			Id:            accountID.String(),
			Acc:           &pb.Account{UserId: userID.String(), Name: "Updated Name"},
			DeleteFileIds: []string{uuid.New().String()},
		}
		existingAcc := &domain.Account{ID: accountID, UserID: userID}

		repo.On("GetAccountByID", ctx, accountID).Return(existingAcc, nil).Once()

		assetHelper.On("SyncFiles", ctx, mock.MatchedBy(func(p domain.FileSyncParams) bool {
			return p.EntityID == accountID && p.EntityType == "account"
		})).Return(nil).Once()

		repo.On("UpdateAccount", ctx, mock.Anything).Return(existingAcc, nil).Once()

		res, err := uc.UpdateAccount(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
		assetHelper.AssertExpectations(t)
	})

	t.Run("DeleteAccount_Success", func(t *testing.T) {
		req := &pb.DeleteAssetRequest{Id: accountID.String(), UserId: userID.String()}

		repo.On("GetAccountByID", ctx, accountID).Return(&domain.Account{ID: accountID}, nil).Once()
		repo.On("SoftDeleteAccount", ctx, accountID, userID).Return(nil).Once()

		res, err := uc.DeleteAccount(ctx, req)
		assert.NoError(t, err)
		assert.True(t, res.Success)
	})

	t.Run("CleanupExpiredAccounts_Success", func(t *testing.T) {
		expiredAccs := []domain.Account{
			{ID: accountID, Files: []domain.FileAssociate{{ID: uuid.New()}}},
		}

		repo.On("GetExpiredAccounts", ctx, mock.AnythingOfType("time.Time")).Return(expiredAccs, nil).Once()

		assetHelper.On("CleanupResource", ctx, accountID, expiredAccs[0].Files, mock.Anything).
			Run(func(args mock.Arguments) {
				callback := args.Get(3).(func(uuid.UUID) error)
				callback(accountID)
			}).Return().Once()

		repo.On("HardDeleteAccount", ctx, accountID).Return(nil).Once()

		err := uc.CleanupExpiredAccounts(ctx)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
		assetHelper.AssertExpectations(t)
	})

	t.Run("UpdateAccount_InvalidID", func(t *testing.T) {
		req := &pb.UpdateAccountRequest{Id: "invalid-uuid"}
		res, err := uc.UpdateAccount(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
