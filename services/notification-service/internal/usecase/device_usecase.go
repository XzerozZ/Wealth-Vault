package usecase

import (
	"context"
	"time"
	"wealth-vault/notification-service/internal/domain"
	repo "wealth-vault/notification-service/internal/repository/interface"

	"github.com/google/uuid"
)

type DeviceUsecase struct {
	repo repo.DeviceRepository
}

func NewDeviceUsecase(repo repo.DeviceRepository) *DeviceUsecase {
	return &DeviceUsecase{
		repo: repo,
	}
}

func (u *DeviceUsecase) RegisterDevice(ctx context.Context, userID uuid.UUID, req *domain.RegisterDeviceRequest) error {
	token := domain.DeviceToken{
		ID:         uuid.New(),
		UserID:     userID,
		Token:      req.Token,
		Platform:   req.Platform,
		DeviceName: req.DeviceName,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return u.repo.RegisterDevice(ctx, &token)
}

func (u *DeviceUsecase) UnregisterDevice(ctx context.Context, userID uuid.UUID, token string) error {
	return u.repo.UnregisterDevice(ctx, userID, token)
}

func (u *DeviceUsecase) GetDevices(ctx context.Context, userID uuid.UUID) ([]domain.DeviceToken, error) {
	return u.repo.GetActiveTokens(ctx, userID)
}
