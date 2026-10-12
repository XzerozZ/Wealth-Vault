package line

import "context"

type LineClient interface {
	SendTextMessage(ctx context.Context, to string, text string, retryKey string) error
}
