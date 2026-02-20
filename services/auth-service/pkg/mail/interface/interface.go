package mail

import (
	"context"
	"wealth-vault/auth-service/internal/domain"
)

type NotificationClient interface {
	SendOTP(ctx context.Context, req domain.SendEmailRequest) error
}
