package push_provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	a "wealth-vault/notification-service/internal/infra/push_provider/interface"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FCMProvider struct {
	client *messaging.Client
}

func NewFCMProvider(base64Credentials string) (*FCMProvider, error) {
	ctx := context.Background()
	decodedJSON, err := base64.StdEncoding.DecodeString(base64Credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 FCM credentials: %w", err)
	}

	opt := option.WithCredentialsJSON(decodedJSON)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Messaging client: %v", err)
	}

	return &FCMProvider{
		client: client,
	}, nil
}

func (f *FCMProvider) Send(ctx context.Context, tokens []string, payload a.PushPayload) []a.PushResult {
	if len(tokens) == 0 {
		return nil
	}

	results := make([]a.PushResult, 0, len(tokens))
	for i := 0; i < len(tokens); i += 500 {
		end := i + 500
		if end > len(tokens) {
			end = len(tokens)
		}
		results = append(results, f.sendBatch(ctx, tokens[i:end], payload)...)
	}
	return results
}

func (f *FCMProvider) sendBatch(ctx context.Context, tokens []string, payload a.PushPayload) []a.PushResult {
	msg := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title:    payload.Title,
			Body:     payload.Body,
			ImageURL: payload.ImageURL,
		},
		Data: payload.Data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				ChannelID: "default",
				Sound:     "default",
			},
		},
		Webpush: &messaging.WebpushConfig{
			Notification: &messaging.WebpushNotification{
				Title: payload.Title,
				Body:  payload.Body,
				Icon:  payload.ImageURL,
			},
		},
	}

	batchResp, err := f.client.SendEachForMulticast(ctx, msg)
	if err != nil {
		results := make([]a.PushResult, len(tokens))
		for i, t := range tokens {
			results[i] = a.PushResult{Token: t, Success: false, Err: err}
		}

		return results
	}

	results := make([]a.PushResult, len(tokens))
	for i, resp := range batchResp.Responses {
		results[i] = a.PushResult{Token: tokens[i], Success: resp.Success}
		if !resp.Success {
			results[i].Err = resp.Error
			results[i].Invalid = isFCMTokenInvalid(resp.Error)
			log.Printf("FCM: failed token %.20s...: %v", tokens[i], resp.Error)
		}
	}

	log.Printf("FCM: %d/%d sent successfully", batchResp.SuccessCount, len(tokens))
	return results
}

func isFCMTokenInvalid(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()
	return strings.Contains(msg, "registration-token-not-registered") ||
		strings.Contains(msg, "invalid-registration-token") ||
		strings.Contains(msg, "token has been deleted")
}
