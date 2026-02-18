package usecase

import (
	"context"
	"wealth-vault/notification-service/internal/domain"

	"github.com/google/uuid"
)

type NotificationUsecase interface {
	HandleGroupCreated(ctx context.Context, evt domain.GroupCreatedEvent) error
	HandleGroupMemberAdded(ctx context.Context, evt domain.GroupMemberAddedEvent) error
	HandleItemShared(ctx context.Context, evt domain.ItemSharedEvent) error
	HandleFriendRequest(ctx context.Context, evt domain.FriendRequestEvent) error
	HandleFriendAccepted(ctx context.Context, evt domain.FriendAcceptedEvent) error
	HandleAccessGranted(ctx context.Context, evt domain.AccessGrantedEvent) error
	HandleMemberRemoved(ctx context.Context, evt domain.MemberRemovedEvent) error
	HandleGroupActivity(ctx context.Context, evt domain.GroupActivityEvent) error
	HandleInsuranceExpiring(ctx context.Context, evt domain.InsuranceExpiringEvent) error
	GetHistory(ctx context.Context, userID uuid.UUID) ([]domain.Notification, error)
}
