package line

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"wealth-vault/notification-service/internal/domain"
)

type lineClient struct {
	channelAccessToken string
	client             *http.Client
}

func NewLineClient(token string) LineClient {
	return &lineClient{
		channelAccessToken: token,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *lineClient) SendTextMessage(to string, text string) error {

	url := "https://api.line.me/v2/bot/message/push"

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

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.channelAccessToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("line request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("line api error status=%d body=%s", resp.StatusCode, string(body))
	}

	return nil
}
