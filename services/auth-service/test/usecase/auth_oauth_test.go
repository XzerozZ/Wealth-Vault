package usecase_test

import (
	"context"
	"errors"
	"testing"
	"wealth-vault/auth-service/internal/domain"
	"wealth-vault/auth-service/internal/usecase"
	pb "wealth-vault/auth-service/pkg/pb/proto/user"
	mock_client "wealth-vault/auth-service/test/mock/client"
	mock_google "wealth-vault/auth-service/test/mock/pkg/google"
	mock_token "wealth-vault/auth-service/test/mock/pkg/token"
	mock_repo "wealth-vault/auth-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/api/idtoken"
)

func TestLoginWithGoogle(t *testing.T) {
	googleIDToken := "valid_google_token"
	email := "googleuser@example.com"
	subjectID := "google_sub_123456789"
	userID := uuid.New()

	validPayload := &idtoken.Payload{
		Claims: map[string]interface{}{
			"email": email,
			"name":  "Google User",
		},
		Subject: subjectID,
	}

	t.Run("success - existing user (Login)", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		mockToken := new(mock_token.MockTokenMaker)
		mockValidator := new(mock_google.MockGoogleValidator)

		uc := usecase.NewAuthUsecase(mockRepo, nil, mockToken, nil, mockValidator)

		mockValidator.
			On("Validate", mock.Anything, googleIDToken).
			Return(validPayload, nil)

		mockRepo.
			On("FindByEmailAndProvider", mock.Anything, email, usecase.ProviderGoogle).
			Return(&domain.AuthAccount{
				UserID:   userID,
				Email:    email,
				Provider: usecase.ProviderGoogle,
			}, nil)

		mockToken.
			On("CreateToken", userID.String(), email, usecase.TokenTypeAccess, mock.Anything).
			Return("mock_access_token", nil)
		mockToken.
			On("CreateToken", userID.String(), email, usecase.TokenTypeRefresh, mock.Anything).
			Return("mock_refresh_token", nil)
		mockRepo.
			On("CreateSession", mock.Anything, mock.AnythingOfType("*domain.AuthSession")).
			Return(nil)

		res, err := uc.LoginWithGoogle(context.Background(), googleIDToken)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, userID.String(), res.UserID)
		assert.Equal(t, "mock_access_token", res.AccessToken)

		mockValidator.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockToken.AssertExpectations(t)
	})

	t.Run("success - new user (Auto Register)", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		mockUserClient := new(mock_client.MockUserClient)
		mockToken := new(mock_token.MockTokenMaker)
		mockValidator := new(mock_google.MockGoogleValidator)

		uc := usecase.NewAuthUsecase(mockRepo, mockUserClient, mockToken, nil, mockValidator)

		mockValidator.
			On("Validate", mock.Anything, googleIDToken).
			Return(validPayload, nil)

		mockRepo.
			On("FindByEmailAndProvider", mock.Anything, email, usecase.ProviderGoogle).
			Return((*domain.AuthAccount)(nil), nil)

		newUserID := uuid.New()
		mockUserClient.
			On("CreateUser", mock.Anything, &pb.CreateUserRequest{
				Email:    email,
				Username: "Google User",
			}).
			Return(&pb.UserResponse{
				User: &pb.User{Id: newUserID.String()},
			}, nil)

		mockRepo.
			On("Register", mock.Anything, mock.AnythingOfType("*domain.AuthAccount")).
			Return(nil)

		mockToken.
			On("CreateToken", newUserID.String(), email, usecase.TokenTypeAccess, mock.Anything).
			Return("mock_access_token", nil)
		mockToken.
			On("CreateToken", newUserID.String(), email, usecase.TokenTypeRefresh, mock.Anything).
			Return("mock_refresh_token", nil)
		mockRepo.
			On("CreateSession", mock.Anything, mock.AnythingOfType("*domain.AuthSession")).
			Return(nil)

		res, err := uc.LoginWithGoogle(context.Background(), googleIDToken)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, newUserID.String(), res.UserID)

		mockValidator.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockUserClient.AssertExpectations(t)
		mockToken.AssertExpectations(t)
	})

	t.Run("error - invalid token", func(t *testing.T) {
		mockValidator := new(mock_google.MockGoogleValidator)

		uc := usecase.NewAuthUsecase(nil, nil, nil, nil, mockValidator)

		mockValidator.
			On("Validate", mock.Anything, "invalid_token").
			Return((*idtoken.Payload)(nil), errors.New("token expired"))

		res, err := uc.LoginWithGoogle(context.Background(), "invalid_token")

		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "invalid google id token")
	})

	t.Run("error - missing email in token", func(t *testing.T) {
		mockValidator := new(mock_google.MockGoogleValidator)

		uc := usecase.NewAuthUsecase(nil, nil, nil, nil, mockValidator)

		invalidPayload := &idtoken.Payload{
			Claims: map[string]interface{}{
				"name": "No Email User",
			},
		}

		mockValidator.
			On("Validate", mock.Anything, googleIDToken).
			Return(invalidPayload, nil)

		res, err := uc.LoginWithGoogle(context.Background(), googleIDToken)

		assert.Nil(t, res)
		assert.Equal(t, "email not found in google token", err.Error())
	})
}
