package mail_test

import (
	"context"
	"testing"

	"wealth-vault/auth-service/configs"
	"wealth-vault/auth-service/internal/domain"
	"wealth-vault/auth-service/pkg/mail"

	"github.com/stretchr/testify/assert"
)

func TestNewMailClient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := configs.Mail{
			Host:   "smtp.example.com",
			Port:   "587",
			Sender: "noreply@wealthvault.com",
			Key:    "secret",
		}

		client, err := mail.NewMailClient(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("error - invalid port", func(t *testing.T) {
		cfg := configs.Mail{
			Host: "smtp.example.com",
			Port: "invalid",
		}

		client, err := mail.NewMailClient(cfg)

		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestSendOTP(t *testing.T) {
	t.Run("dial network error with default expiry", func(t *testing.T) {
		cfg := configs.Mail{
			Host:   "fake.smtp.local",
			Port:   "2525",
			Sender: "test@example.com",
			Key:    "secret",
		}

		client, err := mail.NewMailClient(cfg)
		assert.NoError(t, err)

		req := domain.SendEmailRequest{
			ToEmail: "target@example.com",
			OTP:     "112233",
		}

		err = client.SendOTP(context.Background(), req)

		assert.Error(t, err)
	})

	t.Run("dial network error with custom expiry", func(t *testing.T) {
		cfg := configs.Mail{
			Host:   "fake.smtp.local",
			Port:   "2525",
			Sender: "test@example.com",
			Key:    "secret",
		}

		client, err := mail.NewMailClient(cfg)
		assert.NoError(t, err)

		req := domain.SendEmailRequest{
			ToEmail:   "target@example.com",
			OTP:       "998877",
			ExpiredAt: "15 นาที",
		}

		err = client.SendOTP(context.Background(), req)

		assert.Error(t, err)
	})
}
