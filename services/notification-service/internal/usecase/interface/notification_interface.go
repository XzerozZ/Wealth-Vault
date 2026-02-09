package usecase

import (
	"wealth-vault/notification-service/internal/domain"
)

type NotificationUsecase interface {
	HandleGroupMemberAdded(evt domain.GroupMemberAddedEvent)
	HandleItemShared(evt domain.ItemSharedEvent)
	HandleInsuranceExpiring(evt domain.InsuranceExpiringEvent)
	GetHistory(userIDStr string) ([]domain.Notification, error)
}
