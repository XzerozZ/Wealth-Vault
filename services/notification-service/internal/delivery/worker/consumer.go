package worker

import (
	"encoding/json"
	"log"
	"wealth-vault/notification-service/internal/domain"
	"wealth-vault/notification-service/internal/usecase"

	"github.com/nats-io/nats.go"
)

func StartConsumer(nc *nats.Conn, uc *usecase.NotificationUsecase) {
	log.Println("🎧 NATS Consumer Started: Listening on noti.group.member.added")

	nc.QueueSubscribe("noti.group.created", "noti-workers", func(m *nats.Msg) {
		var evt domain.GroupCreatedEvent
		json.Unmarshal(m.Data, &evt)
		uc.HandleGroupCreated(evt)
	})

	nc.QueueSubscribe("noti.group.member.added", "noti-workers", func(m *nats.Msg) {
		log.Printf("📥 NATS Received: %s", string(m.Data))

		var evt domain.GroupMemberAddedEvent
		if err := json.Unmarshal(m.Data, &evt); err == nil {
			uc.HandleGroupMemberAdded(evt)
		} else {
			log.Printf("❌ JSON Error: %v", err)
		}
	})

	nc.QueueSubscribe("noti.group.member.removed", "noti-workers", func(m *nats.Msg) {
		log.Printf("📥 NATS Received: %s", string(m.Data))

		var evt domain.MemberRemovedEvent
		if err := json.Unmarshal(m.Data, &evt); err == nil {
			uc.HandleMemberRemoved(evt)
		} else {
			log.Printf("❌ JSON Error: %v", err)
		}
	})

	nc.QueueSubscribe("noti.access.granted", "noti-workers", func(m *nats.Msg) {
		log.Printf("📥 NATS Received: %s", string(m.Data))

		var evt domain.AccessGrantedEvent
		if err := json.Unmarshal(m.Data, &evt); err == nil {
			uc.HandleAccessGranted(evt)
		} else {
			log.Printf("❌ JSON Error: %v", err)
		}
	})

	nc.QueueSubscribe("noti.item.shared", "noti-workers", func(m *nats.Msg) {
		log.Printf("📥 NATS Received (Shared): %s", string(m.Data))

		var evt domain.ItemSharedEvent
		if err := json.Unmarshal(m.Data, &evt); err == nil {
			uc.HandleItemShared(evt)
		} else {
			log.Printf("❌ JSON Unmarshal Error: %v", err)
		}
	})

	nc.QueueSubscribe("noti.insurance.expiring", "noti-workers", func(m *nats.Msg) {
		var evt domain.InsuranceExpiringEvent
		if err := json.Unmarshal(m.Data, &evt); err == nil {
			uc.HandleInsuranceExpiring(evt)
		}
	})

	nc.QueueSubscribe("noti.friend.request", "noti-workers", func(m *nats.Msg) {
		log.Printf("📥 NATS Received (Friend Request): %s", string(m.Data))

		var evt domain.FriendRequestEvent
		if err := json.Unmarshal(m.Data, &evt); err == nil {
			uc.HandleFriendRequest(evt)
		} else {
			log.Printf("❌ JSON Error: %v", err)
		}
	})

	nc.QueueSubscribe("noti.friend.accepted", "noti-workers", func(m *nats.Msg) {
		log.Printf("📥 NATS Received (Friend Accepted): %s", string(m.Data))

		var evt domain.FriendAcceptedEvent
		if err := json.Unmarshal(m.Data, &evt); err == nil {
			uc.HandleFriendAccepted(evt)
		} else {
			log.Printf("❌ JSON Error: %v", err)
		}
	})

	nc.QueueSubscribe("noti.group.activity", "noti-workers", func(m *nats.Msg) {
		log.Printf("📥 NATS Received (Group Activity): %s", string(m.Data))

		var evt domain.GroupActivityEvent
		if err := json.Unmarshal(m.Data, &evt); err == nil {
			uc.HandleGroupActivity(evt)
		} else {
			log.Printf("❌ JSON Error: %v", err)
		}
	})
}
