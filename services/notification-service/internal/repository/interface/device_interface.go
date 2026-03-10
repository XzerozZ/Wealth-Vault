package reepository

import (
	"context"
	"wealth-vault/notification-service/internal/domain"

	"github.com/google/uuid"
)

type DeviceRepository interface {
	RegisterDevice(ctx context.Context, req *domain.DeviceToken) error
	GetActiveTokens(ctx context.Context, userID uuid.UUID) ([]domain.DeviceToken, error)
	UnregisterDevice(ctx context.Context, userID uuid.UUID, token string) error
	MarkTokenInactive(ctx context.Context, token string) error
}
