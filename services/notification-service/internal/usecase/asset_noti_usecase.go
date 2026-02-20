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

	message := m.BuildInsuranceExpireMessage(evt)

	noti := &domain.Notification{
		ID:         uuid.New(),
		EntityType: "INSURANCE",
		EntityID:   insuranceID,
		Receiver:   userID,
		Message:    message,
		Channel:    "IN_APP",
		CreatedAt:  time.Now(),
		IsRead:     false,
	}

	if err := u.repo.CreateNotification(ctx, noti); err != nil {
		return fmt.Errorf("create notification: %w", err)
	}

	u.hub.Emit(userID.String(), noti)

	return nil
}
