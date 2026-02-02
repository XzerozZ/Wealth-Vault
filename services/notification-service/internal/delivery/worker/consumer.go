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

	nc.QueueSubscribe("noti.group.member.added", "noti-workers", func(m *nats.Msg) {
		log.Printf("📥 NATS Received: %s", string(m.Data))

		var evt domain.GroupMemberAddedEvent
		if err := json.Unmarshal(m.Data, &evt); err == nil {
			uc.HandleGroupMemberAdded(evt)
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
}
