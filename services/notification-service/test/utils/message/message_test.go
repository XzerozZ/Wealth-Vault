package message_test

import (
	"testing"
	"wealth-vault/notification-service/internal/domain"
	message "wealth-vault/notification-service/pkg/utils/message"

	"github.com/stretchr/testify/assert"
)

func TestBuildInsuranceExpireMessage(t *testing.T) {
	insName := "ประกันสุขภาพ AIA"
	expDate := "20/02/2026"

	tests := []struct {
		name     string
		input    domain.InsuranceExpiringEvent
		expected string
	}{
		{
			name: "Success - Urgent Message (1 day left)",
			input: domain.InsuranceExpiringEvent{
				InsuranceName: insName,
				DaysLeft:      1,
				ExpDate:       expDate,
			},
			expected: "⚠️ ด่วน! ประกัน 'ประกันสุขภาพ AIA' จะหมดอายุในวันพรุ่งนี้ (20/02/2026)",
		},
		{
			name: "Success - Weekly Message (7 days left)",
			input: domain.InsuranceExpiringEvent{
				InsuranceName: insName,
				DaysLeft:      7,
				ExpDate:       expDate,
			},
			expected: "📢 แจ้งเตือน: ประกัน 'ประกันสุขภาพ AIA' จะหมดอายุในอีก 1 อาทิตย์ (20/02/2026)",
		},
		{
			name: "Success - General Warning (e.g. 3 days left)",
			input: domain.InsuranceExpiringEvent{
				InsuranceName: insName,
				DaysLeft:      3,
				ExpDate:       expDate,
			},
			expected: "📢 แจ้งเตือน: ประกัน 'ประกันสุขภาพ AIA' จะหมดอายุในอีก 3 วัน (20/02/2026)",
		},
		{
			name: "Success - Far Warning (e.g. 30 days left)",
			input: domain.InsuranceExpiringEvent{
				InsuranceName: insName,
				DaysLeft:      30,
				ExpDate:       expDate,
			},
			expected: "📢 แจ้งเตือน: ประกัน 'ประกันสุขภาพ AIA' จะหมดอายุในอีก 30 วัน (20/02/2026)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := message.BuildInsuranceExpireMessage(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
