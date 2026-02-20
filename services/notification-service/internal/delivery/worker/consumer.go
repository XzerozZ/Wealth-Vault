package worker

import (
	"context"
	"log"
	"time"
	"wealth-vault/notification-service/internal/domain"
	"wealth-vault/notification-service/internal/usecase"

	"github.com/nats-io/nats.go"
)

const (
	queueGroup     = "noti-workers"
	handlerTimeout = 10 * time.Second
)

func StartConsumer(nc *nats.Conn, uc *usecase.NotificationUsecase) error {
	log.Println("🎧 NATS Consumer Started")

	subscriptions := []struct {
		subject string
		handler func(context.Context, []byte) error
	}{
		{"noti.group.created", wrap[domain.GroupCreatedEvent](uc.HandleGroupCreated)},
		{"noti.group.member.added", wrap[domain.GroupMemberAddedEvent](uc.HandleGroupMemberAdded)},
		{"noti.group.member.removed", wrap[domain.MemberRemovedEvent](uc.HandleMemberRemoved)},
		{"noti.access.granted", wrap[domain.AccessGrantedEvent](uc.HandleAccessGranted)},
		{"noti.item.shared", wrap[domain.ItemSharedEvent](uc.HandleItemShared)},
		{"noti.friend.request", wrap[domain.FriendRequestEvent](uc.HandleFriendRequest)},
		{"noti.friend.accepted", wrap[domain.FriendAcceptedEvent](uc.HandleFriendAccepted)},
		{"noti.group.activity", wrap[domain.GroupActivityEvent](uc.HandleGroupActivity)},
		{"noti.insurance.expiring", wrap[domain.InsuranceExpiringEvent](uc.HandleInsuranceExpiring)},
	}

	for _, sub := range subscriptions {
		sub := sub

		_, err := nc.QueueSubscribe(sub.subject, queueGroup, func(m *nats.Msg) {
			ctx, cancel := context.WithTimeout(context.Background(), handlerTimeout)
			defer cancel()

			log.Printf("📥 NATS Received [%s]", sub.subject)

			if err := sub.handler(ctx, m.Data); err != nil {
				log.Printf("❌ Handler Error [%s]: %v", sub.subject, err)
			}
		})
		if err != nil {
			return err
		}
	}

	return nil
}
