package usecase

import (
	"fmt"
	"log"
	"time"
	"wealth-vault/notification-service/internal/domain"

	"github.com/google/uuid"
)

func (u *NotificationUsecase) HandleInsuranceExpiring(evt domain.InsuranceExpiringEvent) {
	userID, _ := uuid.Parse(evt.UserID)
	insuranceID, _ := uuid.Parse(evt.InsuranceID)

	var message string
	if evt.DaysLeft == 1 {
		message = fmt.Sprintf("⚠️ ด่วน! ประกัน '%s' จะหมดอายุในวันพรุ่งนี้ (%s)", evt.InsuranceName, evt.ExpDate)
	} else {
		message = fmt.Sprintf("📢 แจ้งเตือน: ประกัน '%s' จะหมดอายุในอีก 1 อาทิตย์ (%s)", evt.InsuranceName, evt.ExpDate)
	}

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

	if err := u.repo.CreateNotification(noti); err != nil {
		log.Printf("❌ Save DB Error: %v", err)
	}

	u.hub.Emit(evt.UserID, noti)
}
