package domain

import "time"

type Message struct {
	ID          string                 `json:"id"`
	SenderID    string                 `json:"sender_id"`
	MsgType     string                 `json:"msg_type"`
	Content     string                 `json:"content"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	SenderName  string                 `json:"sender_name"`
	SenderImage string                 `json:"sender_image"`
	IsMe        bool                   `json:"is_me"`
}
