package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"wealth-vault/auth-service/internal/domain"
	"wealth-vault/auth-service/internal/usecase"
	authPb "wealth-vault/auth-service/pkg/pb/proto/auth"
	mock_mail "wealth-vault/auth-service/test/mock/pkg/mail"
	mock_token "wealth-vault/auth-service/test/mock/pkg/token"
	mock_repo "wealth-vault/auth-service/test/mock/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestForgotPassword(t *testing.T) {
	email := "test@example.com"
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		mockMail := new(mock_mail.MockNotificationClient)

		uc := usecase.NewAuthUsecase(mockRepo, nil, nil, mockMail, nil)

		req := &authPb.ForgotPasswordRequest{
			Email: email,
		}

		mockRepo.
			On("FindByEmailAndProvider", mock.Anything, email, usecase.ProviderLocal).
			Return(&domain.AuthAccount{
				UserID:   userID,
				Email:    email,
				Provider: usecase.ProviderLocal,
			}, nil)

		mockRepo.
			On("SaveOTP", mock.Anything, mock.AnythingOfType("*domain.AuthOTP")).
			Return(nil)

		mockMail.
			On("SendOTP", mock.Anything, mock.AnythingOfType("domain.SendEmailRequest")).
			Return(nil)

		res, err := uc.ForgotPassword(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.True(t, res.Success)

		time.Sleep(50 * time.Millisecond)

		mockRepo.AssertExpectations(t)
		mockMail.AssertExpectations(t)
	})

	t.Run("error - user not found", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)

		uc := usecase.NewAuthUsecase(mockRepo, nil, nil, nil, nil)

		mockRepo.
			On("FindByEmailAndProvider", mock.Anything, "wrong@example.com", usecase.ProviderLocal).
			Return((*domain.AuthAccount)(nil), errors.New("record not found"))

		res, err := uc.ForgotPassword(context.Background(), &authPb.ForgotPasswordRequest{
			Email: "wrong@example.com",
		})

		assert.Nil(t, res)
		assert.Equal(t, "user not found", err.Error())
	})

	t.Run("error - try to reset social login account", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		uc := usecase.NewAuthUsecase(mockRepo, nil, nil, nil, nil)

		mockRepo.
			On("FindByEmailAndProvider", mock.Anything, email, usecase.ProviderLocal).
			Return(&domain.AuthAccount{
				UserID:   userID,
				Email:    email,
				Provider: usecase.ProviderGoogle,
			}, nil)

		res, err := uc.ForgotPassword(context.Background(), &authPb.ForgotPasswordRequest{
			Email: email,
		})

		assert.Nil(t, res)
		assert.Equal(t, "cannot reset password for social login account", err.Error())
	})
}

func TestVerifyForgotPasswordOTP(t *testing.T) {
	email := "test@example.com"
	userID := uuid.New()
	accountID := uuid.New()
	otpCode := "123456"

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		mockToken := new(mock_token.MockTokenMaker)

		uc := usecase.NewAuthUsecase(mockRepo, nil, mockToken, nil, nil)

		req := &authPb.VerifyOTPRequest{
			Email: email,
			Otp:   otpCode,
		}

		mockRepo.
			On("FindByEmailAndProvider", mock.Anything, email, usecase.ProviderLocal).
			Return(&domain.AuthAccount{
				ID:       accountID,
				UserID:   userID,
				Email:    email,
				Provider: usecase.ProviderLocal,
			}, nil)

		mockRepo.
			On("GetValidOTP", mock.Anything, userID, otpCode).
			Return(&domain.AuthOTP{OTP: otpCode}, nil)

		mockToken.
			On("CreateToken", userID.String(), email, usecase.TokenTypeReset, mock.Anything).
			Return("mock_reset_token", nil)

		mockRepo.
			On("DeleteOTP", mock.Anything, accountID).
			Return(nil)

		res, err := uc.VerifyForgotPasswordOTP(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.True(t, res.Success)
		assert.Equal(t, "mock_reset_token", res.ResetToken)

		mockRepo.AssertExpectations(t)
		mockToken.AssertExpectations(t)
	})

	t.Run("error - invalid OTP", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)

		uc := usecase.NewAuthUsecase(mockRepo, nil, nil, nil, nil)

		mockRepo.
			On("FindByEmailAndProvider", mock.Anything, email, usecase.ProviderLocal).
			Return(&domain.AuthAccount{UserID: userID, Email: email}, nil)

		mockRepo.
			On("GetValidOTP", mock.Anything, userID, "wrong_otp").
			Return((*domain.AuthOTP)(nil), errors.New("otp expired"))

		res, err := uc.VerifyForgotPasswordOTP(context.Background(), &authPb.VerifyOTPRequest{
			Email: email,
			Otp:   "wrong_otp",
		})

		assert.Nil(t, res)
		assert.Equal(t, "invalid or expired OTP", err.Error())
	})
}

func TestResetPassword(t *testing.T) {
	email := "test@example.com"
	userID := uuid.New()
	oldPassword := "old_password_123"
	newPassword := "new_password_456"

	hashedOldPwd, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		mockToken := new(mock_token.MockTokenMaker)

		uc := usecase.NewAuthUsecase(mockRepo, nil, mockToken, nil, nil)

		req := &authPb.ResetPasswordRequest{
			ResetToken:  "valid_reset_token",
			NewPassword: newPassword,
		}

		mockToken.
			On("VerifyToken", "valid_reset_token").
			Return(jwt.MapClaims{
				"type":  usecase.TokenTypeReset,
				"email": email,
			}, nil)

		mockRepo.
			On("FindByEmailAndProvider", mock.Anything, email, usecase.ProviderLocal).
			Return(&domain.AuthAccount{
				UserID:   userID,
				Email:    email,
				Password: string(hashedOldPwd),
			}, nil)

		mockRepo.
			On("UpdatePassword", mock.Anything, userID, mock.AnythingOfType("string")).
			Return(nil)

		res, err := uc.ResetPassword(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.True(t, res.Success)

		mockRepo.AssertExpectations(t)
		mockToken.AssertExpectations(t)
	})

	t.Run("error - same as old password", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)
		mockToken := new(mock_token.MockTokenMaker)

		uc := usecase.NewAuthUsecase(mockRepo, nil, mockToken, nil, nil)

		mockToken.
			On("VerifyToken", "valid_reset_token").
			Return(jwt.MapClaims{
				"type":  usecase.TokenTypeReset,
				"email": email,
			}, nil)

		mockRepo.
			On("FindByEmailAndProvider", mock.Anything, email, usecase.ProviderLocal).
			Return(&domain.AuthAccount{
				UserID:   userID,
				Email:    email,
				Password: string(hashedOldPwd),
			}, nil)

		req := &authPb.ResetPasswordRequest{
			ResetToken:  "valid_reset_token",
			NewPassword: oldPassword,
		}

		res, err := uc.ResetPassword(context.Background(), req)

		assert.Nil(t, res)
		assert.Equal(t, "new password cannot be the same as the old password", err.Error())
	})

	t.Run("error - invalid token type", func(t *testing.T) {
		mockToken := new(mock_token.MockTokenMaker)
		uc := usecase.NewAuthUsecase(nil, nil, mockToken, nil, nil)

		mockToken.
			On("VerifyToken", "access_token").
			Return(jwt.MapClaims{
				"type":  usecase.TokenTypeAccess,
				"email": email,
			}, nil)

		res, err := uc.ResetPassword(context.Background(), &authPb.ResetPasswordRequest{
			ResetToken:  "access_token",
			NewPassword: newPassword,
		})

		assert.Nil(t, res)
		assert.Equal(t, "invalid token type: expected reset token", err.Error())
	})
}

func TestDeleteExpiredOTPs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)

		uc := usecase.NewAuthUsecase(mockRepo, nil, nil, nil, nil)

		mockRepo.On("DeleteExpiredOTPs", mock.Anything).Return(nil)

		err := uc.DeleteExpiredOTPs(context.Background())

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mock_repo.MockAuthRepository)

		uc := usecase.NewAuthUsecase(mockRepo, nil, nil, nil, nil)

		mockRepo.On("DeleteExpiredOTPs", mock.Anything).Return(errors.New("db error"))

		err := uc.DeleteExpiredOTPs(context.Background())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to cleanup OTPs")
	})
}
