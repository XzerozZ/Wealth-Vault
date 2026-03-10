package usecase_test

import (
	"context"
	"errors"
	"testing"
	"wealth-vault/auth-service/internal/domain"
	"wealth-vault/auth-service/internal/usecase"
	mock_repo "wealth-vault/auth-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLinkLineAccount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		uc := usecase.NewAuthUsecase(mockRepo, nil, nil, nil, nil)

		userID := uuid.New()
		lineUserID := "line-123456"
		ctx := context.Background()

		mockRepo.
			On("FindByUserIDAndProvider", ctx, userID, usecase.ProviderLine).
			Return((*domain.AuthAccount)(nil), nil)

		mockRepo.
			On("Register", ctx, mock.MatchedBy(func(auth *domain.AuthAccount) bool {
				return auth.UserID == userID &&
					auth.Provider == usecase.ProviderLine &&
					auth.ProviderAccountID == lineUserID &&
					auth.IsEmailVerified == true
			})).
			Return(nil)

		err := uc.LinkLineAccount(ctx, userID.String(), lineUserID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("already linked", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		uc := usecase.NewAuthUsecase(mockRepo, nil, nil, nil, nil)

		userID := uuid.New()
		lineUserID := "line-123456"

		mockRepo.
			On("FindByUserIDAndProvider", mock.Anything, userID, usecase.ProviderLine).
			Return(&domain.AuthAccount{UserID: userID, Provider: usecase.ProviderLine}, nil)

		err := uc.LinkLineAccount(context.Background(), userID.String(), lineUserID)

		assert.Error(t, err)
		assert.Equal(t, "user already linked with LINE account", err.Error())
		mockRepo.AssertNotCalled(t, "Register", mock.Anything, mock.Anything)
	})

	t.Run("invalid user id format", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		uc := usecase.NewAuthUsecase(mockRepo, nil, nil, nil, nil)

		err := uc.LinkLineAccount(context.Background(), "bad-uuid", "line-123")

		assert.Error(t, err)
		assert.Equal(t, "invalid user id format", err.Error())
	})

	t.Run("register failure", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		uc := usecase.NewAuthUsecase(mockRepo, nil, nil, nil, nil)

		userID := uuid.New()

		mockRepo.On("FindByUserIDAndProvider", mock.Anything, userID, usecase.ProviderLine).
			Return((*domain.AuthAccount)(nil), nil)

		mockRepo.On("Register", mock.Anything, mock.Anything).
			Return(errors.New("db error"))

		err := uc.LinkLineAccount(context.Background(), userID.String(), "line-123")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to link line account")
	})
}
