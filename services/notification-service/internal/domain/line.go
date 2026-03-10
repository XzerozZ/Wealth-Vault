package domain

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type PushMessageRequest struct {
	To       string    `json:"to"`
	Messages []Message `json:"messages"`
}
