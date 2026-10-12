package usecase

import (
	"context"
	"fmt"
	"time"
	"wealth-vault/notification-service/internal/domain"
	m "wealth-vault/notification-service/pkg/utils/message"

	"github.com/google/uuid"
)

func (u *NotificationUsecase) HandleInsuranceExpiring(ctx context.Context, evt domain.InsuranceExpiringEvent) error {
	userID, err := uuid.Parse(evt.UserID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}
	insuranceID, err := uuid.Parse(evt.InsuranceID)
	if err != nil {
		return fmt.Errorf("invalid insurance id: %w", err)
	}
	return u.notifyTarget(ctx, userID, nil, "INSURANCE", insuranceID, m.BuildInsuranceExpireMessage(evt), time.Now().Unix(), nil)
}
