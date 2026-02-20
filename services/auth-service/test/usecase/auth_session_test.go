package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"wealth-vault/auth-service/internal/domain"
	"wealth-vault/auth-service/internal/usecase"
	mock_token "wealth-vault/auth-service/test/mock/pkg/token"
	mock_repo "wealth-vault/auth-service/test/mock/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRefreshToken(t *testing.T) {
	refreshToken := "valid_refresh_token"
	email := "test@example.com"
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		mockToken := new(mock_token.MockTokenMaker)

		uc := usecase.NewAuthUsecase(mockRepo, nil, mockToken, nil, nil)

		mockToken.
			On("VerifyToken", refreshToken).
			Return(jwt.MapClaims{
				"type":  usecase.TokenTypeRefresh,
				"email": email,
			}, nil)

		mockRepo.
			On("GetSessionByRefreshToken", mock.Anything, refreshToken).
			Return(&domain.AuthSession{
				UserID:           userID,
				RefreshToken:     refreshToken,
				RefreshExpiresAt: time.Now().Add(24 * time.Hour),
			}, nil)

		mockRepo.
			On("RevokeSession", mock.Anything, refreshToken).
			Return(nil)

		mockRepo.
			On("FindByID", mock.Anything, userID.String()).
			Return(&domain.AuthAccount{
				UserID: userID,
				Email:  email,
			}, nil)

		mockToken.
			On("CreateToken", userID.String(), email, usecase.TokenTypeAccess, mock.Anything).
			Return("new_access_token", nil)
		mockToken.
			On("CreateToken", userID.String(), email, usecase.TokenTypeRefresh, mock.Anything).
			Return("new_refresh_token", nil)

		mockRepo.
			On("CreateSession", mock.Anything, mock.AnythingOfType("*domain.AuthSession")).
			Return(nil)

		res, err := uc.RefreshToken(context.Background(), refreshToken)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, userID.String(), res.UserID)
		assert.Equal(t, "new_access_token", res.AccessToken)
		assert.Equal(t, "new_refresh_token", res.RefreshToken)

		mockRepo.AssertExpectations(t)
		mockToken.AssertExpectations(t)
	})

	t.Run("error - invalid token", func(t *testing.T) {
		mockToken := new(mock_token.MockTokenMaker)
		uc := usecase.NewAuthUsecase(nil, nil, mockToken, nil, nil)

		mockToken.
			On("VerifyToken", "invalid_token").
			Return(jwt.MapClaims(nil), errors.New("jwt expired or invalid"))

		res, err := uc.RefreshToken(context.Background(), "invalid_token")

		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "invalid refresh token")
	})

	t.Run("error - wrong token type", func(t *testing.T) {
		mockToken := new(mock_token.MockTokenMaker)
		uc := usecase.NewAuthUsecase(nil, nil, mockToken, nil, nil)

		mockToken.
			On("VerifyToken", refreshToken).
			Return(jwt.MapClaims{
				"type": usecase.TokenTypeAccess,
			}, nil)

		res, err := uc.RefreshToken(context.Background(), refreshToken)

		assert.Nil(t, res)
		assert.Equal(t, "invalid token type", err.Error())
	})

	t.Run("error - session not found or revoked", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		mockToken := new(mock_token.MockTokenMaker)
		uc := usecase.NewAuthUsecase(mockRepo, nil, mockToken, nil, nil)

		mockToken.
			On("VerifyToken", refreshToken).
			Return(jwt.MapClaims{"type": usecase.TokenTypeRefresh}, nil)

		mockRepo.
			On("GetSessionByRefreshToken", mock.Anything, refreshToken).
			Return((*domain.AuthSession)(nil), errors.New("record not found"))

		res, err := uc.RefreshToken(context.Background(), refreshToken)

		assert.Nil(t, res)
		assert.Equal(t, "invalid or revoked refresh token", err.Error())
	})

	t.Run("error - refresh token expired in db", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		mockToken := new(mock_token.MockTokenMaker)
		uc := usecase.NewAuthUsecase(mockRepo, nil, mockToken, nil, nil)

		mockToken.
			On("VerifyToken", refreshToken).
			Return(jwt.MapClaims{"type": usecase.TokenTypeRefresh}, nil)

		mockRepo.
			On("GetSessionByRefreshToken", mock.Anything, refreshToken).
			Return(&domain.AuthSession{
				RefreshExpiresAt: time.Now().Add(-1 * time.Hour),
			}, nil)

		res, err := uc.RefreshToken(context.Background(), refreshToken)

		assert.Nil(t, res)
		assert.Equal(t, "refresh token expired", err.Error())
	})
}

func TestCleanupSessions(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		uc := usecase.NewAuthUsecase(mockRepo, nil, nil, nil, nil)

		mockRepo.On("DeleteExpiredSessions", mock.Anything).Return(nil)

		err := uc.CleanupSessions(context.Background())

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		uc := usecase.NewAuthUsecase(mockRepo, nil, nil, nil, nil)

		mockRepo.On("DeleteExpiredSessions", mock.Anything).Return(errors.New("db timeout"))

		err := uc.CleanupSessions(context.Background())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to cleanup sessions")
	})
}
