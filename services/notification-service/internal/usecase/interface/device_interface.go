package usecase

import (
	"context"
	"wealth-vault/notification-service/internal/domain"

	"github.com/google/uuid"
)

type DeviceUsecase interface {
	RegisterDevice(ctx context.Context, userID uuid.UUID, req *domain.RegisterDeviceRequest) error
	UnregisterDevice(ctx context.Context, userID uuid.UUID, token string) error
	GetDevices(ctx context.Context, userID uuid.UUID) ([]domain.DeviceToken, error)
}
