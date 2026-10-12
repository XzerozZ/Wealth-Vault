package line

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"wealth-vault/notification-service/internal/domain"
)

type lineClient struct {
	channelAccessToken string
	client             *http.Client
	endpoint           string
}

func NewLineClient(token string) LineClient {
	return &lineClient{
		channelAccessToken: token,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		endpoint: "https://api.line.me/v2/bot/message/push",
	}
}

func (c *lineClient) SendTextMessage(ctx context.Context, to string, text string, retryKey string) error {

	payload := domain.PushMessageRequest{
		To: to,
		Messages: []domain.Message{
			{
				Type: "text",
				Text: text,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload error: %w", err)
	}

	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("create LINE request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.channelAccessToken)
		req.Header.Set("X-Line-Retry-Key", retryKey)
		resp, err := c.client.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusConflict {
				return nil
			}
			if resp.StatusCode < 500 {
				return fmt.Errorf("LINE push rejected with status %d", resp.StatusCode)
			}
		}
		if attempt == 2 {
			if err != nil {
				return fmt.Errorf("LINE push request failed: %w", err)
			}
			return fmt.Errorf("LINE push failed with status %d", resp.StatusCode)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(1<<attempt) * 300 * time.Millisecond):
		}
	}
	return nil
}
