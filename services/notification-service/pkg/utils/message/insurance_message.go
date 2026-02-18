package massage

import (
	"fmt"
	"wealth-vault/notification-service/internal/domain"
)

func BuildInsuranceExpireMessage(evt domain.InsuranceExpiringEvent) string {
	const (
		msgUrgent  = "⚠️ ด่วน! ประกัน '%s' จะหมดอายุในวันพรุ่งนี้ (%s)"
		msgWarning = "📢 แจ้งเตือน: ประกัน '%s' จะหมดอายุในอีก %d วัน (%s)"
		msgWeekly  = "📢 แจ้งเตือน: ประกัน '%s' จะหมดอายุในอีก 1 อาทิตย์ (%s)"
	)

	switch evt.DaysLeft {
	case 1:
		return fmt.Sprintf(msgUrgent, evt.InsuranceName, evt.ExpDate)
	case 7:
		return fmt.Sprintf(msgWeekly, evt.InsuranceName, evt.ExpDate)
	default:
		return fmt.Sprintf(msgWarning, evt.InsuranceName, evt.DaysLeft, evt.ExpDate)
	}
}
