package usecase_test

import (
	"context"
	"errors"
	"testing"

	"wealth-vault/auth-service/internal/domain"
	"wealth-vault/auth-service/internal/usecase"
	pb "wealth-vault/auth-service/pkg/pb/proto/user"
	mock_client "wealth-vault/auth-service/test/mock/client"
	mock_token "wealth-vault/auth-service/test/mock/pkg/token"
	mock_repo "wealth-vault/auth-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		mockUserClient := new(mock_client.MockUserClient)
		uc := usecase.NewAuthUsecase(mockRepo, mockUserClient, nil, nil, nil)

		input := &domain.RegisterInput{
			Email:    "test@example.com",
			Password: "123456",
		}

		mockRepo.
			On("FindByEmailAndProvider", mock.Anything, input.Email, usecase.ProviderLocal).
			Return((*domain.AuthAccount)(nil), nil)

		userID := uuid.New()

		mockUserClient.
			On("CreateUser", mock.Anything, mock.Anything).
			Return(&pb.UserResponse{
				User: &pb.User{
					Id: userID.String(),
				},
			}, nil)

		mockRepo.
			On("Register", mock.Anything, mock.AnythingOfType("*domain.AuthAccount")).
			Return(nil)

		res, err := uc.Register(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, userID.String(), res.UserID)

		mockRepo.AssertExpectations(t)
		mockUserClient.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		mockUserClient := new(mock_client.MockUserClient)

		uc := usecase.NewAuthUsecase(mockRepo, mockUserClient, nil, nil, nil)

		mockRepo.
			On("FindByEmailAndProvider", mock.Anything, "test@example.com", usecase.ProviderLocal).
			Return(&domain.AuthAccount{}, nil)

		res, err := uc.Register(context.Background(), &domain.RegisterInput{
			Email:    "test@example.com",
			Password: "123456",
		})

		assert.Nil(t, res)
		assert.Equal(t, "email already registered with this provider", err.Error())
	})

	t.Run("missing password", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		mockUserClient := new(mock_client.MockUserClient)

		uc := usecase.NewAuthUsecase(mockRepo, mockUserClient, nil, nil, nil)

		mockRepo.
			On("FindByEmailAndProvider", mock.Anything, "test@example.com", usecase.ProviderLocal).
			Return((*domain.AuthAccount)(nil), nil)

		res, err := uc.Register(context.Background(), &domain.RegisterInput{
			Email: "test@example.com",
		})

		assert.Nil(t, res)
		assert.Equal(t, "password is required for local registration", err.Error())
	})
}

func TestLogin(t *testing.T) {
	email := "test@example.com"
	password := "123456"
	userID := uuid.New()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		mockToken := new(mock_token.MockTokenMaker)

		uc := usecase.NewAuthUsecase(mockRepo, nil, mockToken, nil, nil)

		input := &domain.LoginInput{
			Email:    email,
			Password: password,
		}

		mockRepo.
			On("FindByEmailAndProvider", mock.Anything, input.Email, usecase.ProviderLocal).
			Return(&domain.AuthAccount{
				UserID:   userID,
				Email:    email,
				Password: string(hashedPassword),
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

		res, err := uc.Login(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, userID.String(), res.UserID)
		assert.Equal(t, "mock_access_token", res.AccessToken)
		assert.Equal(t, "mock_refresh_token", res.RefreshToken)

		mockRepo.AssertExpectations(t)
		mockToken.AssertExpectations(t)
	})

	t.Run("invalid email (user not found)", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)

		uc := usecase.NewAuthUsecase(mockRepo, nil, nil, nil, nil)

		mockRepo.
			On("FindByEmailAndProvider", mock.Anything, "wrong@example.com", usecase.ProviderLocal).
			Return((*domain.AuthAccount)(nil), errors.New("record not found"))

		res, err := uc.Login(context.Background(), &domain.LoginInput{
			Email:    "wrong@example.com",
			Password: password,
		})

		assert.Nil(t, res)
		assert.Equal(t, "invalid email or password", err.Error())
	})

	t.Run("invalid password", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)

		uc := usecase.NewAuthUsecase(mockRepo, nil, nil, nil, nil)
		mockRepo.
			On("FindByEmailAndProvider", mock.Anything, email, usecase.ProviderLocal).
			Return(&domain.AuthAccount{
				UserID:   userID,
				Email:    email,
				Password: string(hashedPassword),
			}, nil)

		res, err := uc.Login(context.Background(), &domain.LoginInput{
			Email:    email,
			Password: "wrong_password_123",
		})

		assert.Nil(t, res)
		assert.Equal(t, "invalid email or password", err.Error())
	})
}
