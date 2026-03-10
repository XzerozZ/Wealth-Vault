package push_provider

import "context"

type PushPayload struct {
	Title    string
	Body     string
	Data     map[string]string
	ImageURL string
}

type PushResult struct {
	Token   string
	Success bool
	Err     error
	Invalid bool
}

type PushProvider interface {
	Send(ctx context.Context, tokens []string, payload PushPayload) []PushResult
}
