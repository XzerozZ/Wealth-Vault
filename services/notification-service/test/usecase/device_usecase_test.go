package usecase_test

import (
	"context"
	"errors"
	"testing"
	"wealth-vault/notification-service/internal/domain"
	"wealth-vault/notification-service/internal/usecase"
	mock_repo "wealth-vault/notification-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeviceUsecase(t *testing.T) {
	mockRepo := new(mock_repo.MockDeviceRepository)
	uc := usecase.NewDeviceUsecase(mockRepo)

	ctx := context.Background()
	userID := uuid.New()

	t.Run("RegisterDevice_Success", func(t *testing.T) {
		req := &domain.RegisterDeviceRequest{
			Token:      "fcm-token-123",
			Platform:   "ios",
			DeviceName: "iPhone 15",
		}

		mockRepo.On("RegisterDevice", ctx, mock.MatchedBy(func(token *domain.DeviceToken) bool {
			return token.UserID == userID &&
				token.Token == req.Token &&
				token.Platform == req.Platform &&
				token.IsActive == true
		})).Return(nil).Once()

		err := uc.RegisterDevice(ctx, userID, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UnregisterDevice_Success", func(t *testing.T) {
		token := "fcm-token-123"

		mockRepo.On("UnregisterDevice", ctx, userID, token).
			Return(nil).Once()

		err := uc.UnregisterDevice(ctx, userID, token)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetDevices_Success", func(t *testing.T) {
		mockDevices := []domain.DeviceToken{
			{UserID: userID, Token: "token-1", IsActive: true},
			{UserID: userID, Token: "token-2", IsActive: true},
		}

		mockRepo.On("GetActiveTokens", ctx, userID).
			Return(mockDevices, nil).Once()

		res, err := uc.GetDevices(ctx, userID)

		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, "token-1", res[0].Token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetDevices_Error", func(t *testing.T) {
		mockRepo.On("GetActiveTokens", ctx, userID).
			Return(nil, errors.New("repository error")).Once()

		res, err := uc.GetDevices(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "repository error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
