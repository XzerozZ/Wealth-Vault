package domain

type FriendRequestEvent struct {
	SenderID   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	TargetID   string `json:"target_id"`
	OccurredAt int64  `json:"occurred_at"`
}

type FriendAcceptedEvent struct {
	AccepterID   string `json:"accepter_id"`
	AccepterName string `json:"accepter_name"`
	RequesterID  string `json:"requester_id"`
	OccurredAt   int64  `json:"occurred_at"`
}
