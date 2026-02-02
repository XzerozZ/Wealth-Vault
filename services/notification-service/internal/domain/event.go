package domain

type GroupMemberAddedEvent struct {
	GroupID       string   `json:"group_id"`
	SenderID      string   `json:"sender_id"`
	AddedUserIDs  []string `json:"added_user_ids"`
	TargetUserIDs []string `json:"target_user_ids"`
	OccurredAt    int64    `json:"occurred_at"`
}

type ItemSharedEvent struct {
	AssetID       string   `json:"asset_id"`
	SenderID      string   `json:"sender_id"`
	SenderName    string   `json:"sender_name"`
	TargetUserIDs []string `json:"target_user_ids"`
	OccurredAt    int64    `json:"occurred_at"`
}
