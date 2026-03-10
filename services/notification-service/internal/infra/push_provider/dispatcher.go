package push_provider

import (
	"context"
	"log"
	"wealth-vault/notification-service/internal/domain"
	a "wealth-vault/notification-service/internal/infra/push_provider/interface"
)

type TokenCleaner interface {
	MarkTokenInactive(ctx context.Context, token string) error
}

type Dispatcher struct {
	fcm     a.PushProvider
	apns    a.PushProvider
	cleaner TokenCleaner
}

func NewDispatcher(fcm a.PushProvider, apns a.PushProvider, cleaner TokenCleaner) *Dispatcher {
	return &Dispatcher{
		fcm:     fcm,
		apns:    apns,
		cleaner: cleaner,
	}
}

func (d *Dispatcher) SendToUser(ctx context.Context, tokens []domain.DeviceToken, payload a.PushPayload) {
	if len(tokens) == 0 {
		return
	}

	var fcmTokens, apnsTokens []string
	for _, t := range tokens {
		switch t.Platform {
		case domain.PlatformAndroid:
			fcmTokens = append(fcmTokens, t.Token)
		case domain.PlatformIOS:
			apnsTokens = append(apnsTokens, t.Token)
		}
	}

	resultCh := make(chan []a.PushResult, 2)

	go func() {
		if len(fcmTokens) > 0 && d.fcm != nil {
			resultCh <- d.fcm.Send(ctx, fcmTokens, payload)
		} else {
			resultCh <- nil
		}
	}()

	go func() {
		if len(apnsTokens) > 0 && d.apns != nil {
			resultCh <- d.apns.Send(ctx, apnsTokens, payload)
		} else {
			resultCh <- nil
		}
	}()

	for i := 0; i < 2; i++ {
		if results := <-resultCh; results != nil {
			d.cleanInvalidTokens(ctx, results)
		}
	}
}

func (d *Dispatcher) cleanInvalidTokens(ctx context.Context, results []a.PushResult) {
	for _, r := range results {
		if r.Invalid {
			log.Printf("🗑️  Removing invalid token %.20s...", r.Token)
			if err := d.cleaner.MarkTokenInactive(ctx, r.Token); err != nil {
				log.Printf("❌ MarkTokenInactive error: %v", err)
			}
		}
	}
}
