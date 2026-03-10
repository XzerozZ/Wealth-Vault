package push_provider

import (
	"context"
	"wealth-vault/notification-service/internal/domain"
)

type Dispatcher interface {
	SendToUser(ctx context.Context, tokens []domain.DeviceToken, payload PushPayload)
}
