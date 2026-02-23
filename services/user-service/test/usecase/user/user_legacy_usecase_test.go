package usecase_test

import (
	"context"
	"testing"
	"wealth-vault/user-service/internal/domain"
	usecase "wealth-vault/user-service/internal/usecase/user"
	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"
	mock_client "wealth-vault/user-service/test/mock/client"
	mock_repo "wealth-vault/user-service/test/mock/repository"
	mock_usecase "wealth-vault/user-service/test/mock/usecase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessLegacyAutoShare(t *testing.T) {
	ctx := context.Background()

	t.Run("no eligible users", func(t *testing.T) {
		userRepo := new(mock_repo.MockUserRepository)
		assetClient := new(mock_client.MockAssetClient)
		itemUC := new(mock_usecase.MockItemUsecase)

		uc := usecase.NewUserUsecase(userRepo, itemUC, nil, nil, assetClient)

		userRepo.
			On("GetUsersReadyForAutoShare", ctx).
			Return([]domain.User{}, nil)

		err := uc.ProcessLegacyAutoShare(ctx)

		assert.NoError(t, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("one user without friends", func(t *testing.T) {
		userRepo := new(mock_repo.MockUserRepository)
		assetClient := new(mock_client.MockAssetClient)
		itemUC := new(mock_usecase.MockItemUsecase)

		uc := usecase.NewUserUsecase(userRepo, itemUC, nil, nil, assetClient)
		userID := uuid.New()
		userRepo.
			On("GetUsersReadyForAutoShare", ctx).
			Return([]domain.User{
				{ID: userID},
			}, nil)

		err := uc.ProcessLegacyAutoShare(ctx)

		assert.NoError(t, err)
	})
}

func TestProcessSingleUserLegacy(t *testing.T) {
	ctx := context.Background()

	t.Run("no friends", func(t *testing.T) {
		userRepo := new(mock_repo.MockUserRepository)
		assetClient := new(mock_client.MockAssetClient)
		itemUC := new(mock_usecase.MockItemUsecase)

		uc := usecase.NewUserUsecase(userRepo, itemUC, nil, nil, assetClient)

		user := domain.User{
			ID:      uuid.New(),
			Friends: []domain.User{},
		}

		err := uc.ProcessSingleUserLegacy(ctx, user)

		assert.NoError(t, err)
	})

	t.Run("no assets", func(t *testing.T) {
		userRepo := new(mock_repo.MockUserRepository)
		assetClient := new(mock_client.MockAssetClient)
		itemUC := new(mock_usecase.MockItemUsecase)

		uc := usecase.NewUserUsecase(userRepo, itemUC, nil, nil, assetClient)

		userID := uuid.New()
		friendID := uuid.New()

		user := domain.User{
			ID: userID,
			Friends: []domain.User{
				{ID: friendID},
			},
		}

		assetClient.
			On("GetAllAssetIDs", ctx, mock.Anything).
			Return(&assetPb.GetMyAssetsResponse{}, nil)

		userRepo.
			On("MarkAutoShareTriggered", ctx, userID).
			Return(nil)

		err := uc.ProcessSingleUserLegacy(ctx, user)

		assert.NoError(t, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("with assets share success", func(t *testing.T) {
		userRepo := new(mock_repo.MockUserRepository)
		assetClient := new(mock_client.MockAssetClient)
		itemUC := new(mock_usecase.MockItemUsecase)

		uc := usecase.NewUserUsecase(userRepo, itemUC, nil, nil, assetClient)

		userID := uuid.New()
		friendID := uuid.New()

		user := domain.User{
			ID: userID,
			Friends: []domain.User{
				{ID: friendID},
			},
		}

		assets := &assetPb.GetMyAssetsResponse{
			AccountIds: []string{"acc1"},
			CashIds:    []string{"cash1"},
		}

		assetClient.
			On("GetAllAssetIDs", ctx, mock.Anything).
			Return(assets, nil)

		itemUC.
			On("BatchShareAssets", ctx, mock.Anything).
			Return(nil)

		userRepo.
			On("MarkAutoShareTriggered", ctx, userID).
			Return(nil)

		err := uc.ProcessSingleUserLegacy(ctx, user)

		assert.NoError(t, err)

		itemUC.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})
}
