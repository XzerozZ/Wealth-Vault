package usecase

import (
	"context"
	"fmt"
	"log"
	"time"
	"wealth-vault/notification-service/internal/domain"
	pb "wealth-vault/notification-service/pkg/pb/proto/auth"
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

	receiverStr := userID.String()

	if u.hub.IsOnline(receiverStr) {
		u.hub.Emit(receiverStr, domain.WSMessage{
			Type:    "NOTIFICATION",
			Payload: noti,
		})
	} else {
		log.Printf("Push Notification Successfully")
		go u.sendPush(context.Background(), userID, noti.Message, noti)
	}

	go func() {

		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := u.authClient.GetProviderAccount(bgCtx, &pb.GetProviderAccountRequest{
			UserId: receiverStr,
		})

		if err != nil || res == nil {
			log.Printf("GetProviderAccount failed user=%s err=%v", receiverStr, err)
			return
		}

		var lineUserID string

		for _, acc := range res.Accounts {
			if acc.Provider == "line" && acc.IsLinked {
				lineUserID = acc.ProviderAccountId
				break
			}
		}

		if lineUserID == "" {
			return
		}

		if err := u.lineClient.SendTextMessage(lineUserID, message); err != nil {
			log.Printf("LINE notify fail user=%s line=%s err=%v",
				receiverStr,
				lineUserID,
				err,
			)
			return
		}

		log.Printf("LINE notify success user=%s line=%s",
			receiverStr,
			lineUserID,
		)

	}()

	return nil
}
