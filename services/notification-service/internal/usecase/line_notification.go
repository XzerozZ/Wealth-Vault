package usecase

import (
	"context"
	"log"
	"time"
	"wealth-vault/notification-service/internal/domain"
	pb "wealth-vault/notification-service/pkg/pb/proto/auth"
)

func (u *NotificationUsecase) queueLineNotification(noti *domain.Notification) {
	if !u.LineAvailable() {
		return
	}
	message, eligible := lineMessage(noti.EntityType)
	if !eligible {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		enabled, err := u.repo.GetLineEnabled(ctx, noti.Receiver)
		if err != nil || !enabled {
			return
		}
		accounts, err := u.authClient.GetProviderAccount(ctx, &pb.GetProviderAccountRequest{UserId: noti.Receiver.String()})
		if err != nil || accounts == nil {
			log.Printf("LINE notification identity lookup failed: %v", err)
			return
		}
		lineUserID := ""
		for _, account := range accounts.Accounts {
			if account.Provider == "line" && account.IsLinked {
				lineUserID = account.ProviderAccountId
			}
		}
		if lineUserID == "" {
			return
		}
		if err := u.lineClient.SendTextMessage(ctx, lineUserID, message+"\n"+u.lineOptions.LiffURL, noti.ID.String()); err != nil {
			log.Printf("LINE notification delivery failed: %v", err)
		}
	}()
}

func lineMessage(entityType string) (string, bool) {
	switch entityType {
	case "FRIEND_REQUEST":
		return "คุณมีคำขอเป็นเพื่อนใหม่ใน WealthVault", true
	case "FRIEND_ACCEPTED":
		return "คำขอเป็นเพื่อนของคุณได้รับการตอบรับแล้วใน WealthVault", true
	case "ASSET":
		return "มีทรัพย์สินที่แชร์ให้คุณใน WealthVault", true
	case "GROUP_INVITE":
		return "คุณได้รับคำเชิญเข้ากลุ่มใน WealthVault", true
	case "INSURANCE":
		return "คุณมีแจ้งเตือนเกี่ยวกับประกันใน WealthVault", true
	default:
		return "", false
	}
}
