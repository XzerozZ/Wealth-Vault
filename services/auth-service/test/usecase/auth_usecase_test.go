package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"wealth-vault/auth-service/internal/domain"
	"wealth-vault/auth-service/internal/usecase"
	mock_token "wealth-vault/auth-service/test/mock/pkg/token"
	mock_repo "wealth-vault/auth-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		password := "secret_password_123"

		hashed, err := usecase.HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)

		err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
		assert.NoError(t, err)
	})

	t.Run("error - password too long", func(t *testing.T) {
		longPassword := strings.Repeat("a", 73)

		hashed, err := usecase.HashPassword(longPassword)

		assert.Error(t, err)
		assert.Empty(t, hashed)
		assert.Contains(t, err.Error(), "failed to hash password")
	})
}

func TestGenerateTokensAndSession(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	mockAccessToken := "mock_access_token"
	mockRefreshToken := "mock_refresh_token"

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		mockToken := new(mock_token.MockTokenMaker)

		uc := usecase.NewAuthUsecase(mockRepo, nil, mockToken, nil, nil)

		mockToken.
			On("CreateToken", userID.String(), email, usecase.TokenTypeAccess, usecase.AccessTokenTTL).
			Return(mockAccessToken, nil)

		mockToken.
			On("CreateToken", userID.String(), email, usecase.TokenTypeRefresh, usecase.RefreshTokenTTL).
			Return(mockRefreshToken, nil)

		mockRepo.
			On("CreateSession", mock.Anything, mock.MatchedBy(func(s *domain.AuthSession) bool {
				return s.UserID == userID &&
					s.AccessToken == mockAccessToken &&
					s.RefreshToken == mockRefreshToken &&
					!s.Revoked
			})).
			Return(nil)

		res, err := uc.GenerateTokensAndSession(context.Background(), userID, email)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, userID.String(), res.UserID)
		assert.Equal(t, mockAccessToken, res.AccessToken)
		assert.Equal(t, mockRefreshToken, res.RefreshToken)

		mockRepo.AssertExpectations(t)
		mockToken.AssertExpectations(t)
	})

	t.Run("error - failed to create access token", func(t *testing.T) {
		mockToken := new(mock_token.MockTokenMaker)
		uc := usecase.NewAuthUsecase(nil, nil, mockToken, nil, nil)

		mockToken.
			On("CreateToken", userID.String(), email, usecase.TokenTypeAccess, usecase.AccessTokenTTL).
			Return("", errors.New("jwt access error"))

		res, err := uc.GenerateTokensAndSession(context.Background(), userID, email)

		assert.Nil(t, res)
		assert.Equal(t, "failed to create access token: jwt access error", err.Error())
	})

	t.Run("error - failed to create refresh token", func(t *testing.T) {
		mockToken := new(mock_token.MockTokenMaker)
		uc := usecase.NewAuthUsecase(nil, nil, mockToken, nil, nil)

		mockToken.
			On("CreateToken", userID.String(), email, usecase.TokenTypeAccess, usecase.AccessTokenTTL).
			Return(mockAccessToken, nil)

		mockToken.
			On("CreateToken", userID.String(), email, usecase.TokenTypeRefresh, usecase.RefreshTokenTTL).
			Return("", errors.New("jwt refresh error"))

		res, err := uc.GenerateTokensAndSession(context.Background(), userID, email)

		assert.Nil(t, res)
		assert.Equal(t, "failed to create refresh token: jwt refresh error", err.Error())
	})

	t.Run("error - failed to create session in db", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		mockToken := new(mock_token.MockTokenMaker)
		uc := usecase.NewAuthUsecase(mockRepo, nil, mockToken, nil, nil)

		mockToken.
			On("CreateToken", userID.String(), email, usecase.TokenTypeAccess, usecase.AccessTokenTTL).
			Return(mockAccessToken, nil)
		mockToken.
			On("CreateToken", userID.String(), email, usecase.TokenTypeRefresh, usecase.RefreshTokenTTL).
			Return(mockRefreshToken, nil)

		mockRepo.
			On("CreateSession", mock.Anything, mock.AnythingOfType("*domain.AuthSession")).
			Return(errors.New("db connection lost"))

		res, err := uc.GenerateTokensAndSession(context.Background(), userID, email)

		assert.Nil(t, res)
		assert.Equal(t, "failed to create session: db connection lost", err.Error())
	})
}
