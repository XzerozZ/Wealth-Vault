package domain

type GroupCreatedEvent struct {
	GroupID       string   `json:"group_id"`
	GroupName     string   `json:"group_name"`
	SenderID      string   `json:"sender_id"`
	SenderName    string   `json:"sender_name"`
	TargetUserIDs []string `json:"target_user_ids"`
	OccurredAt    int64    `json:"occurred_at"`
}

type GroupActivityEvent struct {
	GroupID      string `json:"group_id"`
	ActivityType string `json:"activity_type"`
	Payload      string `json:"payload"`
	ActorID      string `json:"actor_id"`
	TargetID     string `json:"target_id"`
	OccurredAt   int64  `json:"occurred_at"`
}

type GroupMemberAddedEvent struct {
	GroupID       string   `json:"group_id"`
	SenderID      string   `json:"sender_id"`
	AddedUserIDs  []string `json:"added_user_ids"`
	TargetUserIDs []string `json:"target_user_ids"`
	OccurredAt    int64    `json:"occurred_at"`
}

type MemberRemovedEvent struct {
	GroupID    string `json:"group_id"`
	GroupName  string `json:"group_name"`
	TargetID   string `json:"target_id"`
	ActionBy   string `json:"action_by"`
	OccurredAt int64  `json:"occurred_at"`
}

type ItemSharedEvent struct {
	AssetID       string   `json:"asset_id"`
	SenderID      string   `json:"sender_id"`
	SenderName    string   `json:"sender_name"`
	TargetUserIDs []string `json:"target_user_ids"`
	OccurredAt    int64    `json:"occurred_at"`
}

type AccessGrantedEvent struct {
	GroupID      string `json:"group_id"`
	GrantorID    string `json:"grantor_id"`
	GrantorName  string `json:"grantor_name"`
	TargetUserID string `json:"target_user_id"`
	ItemCount    int    `json:"item_count"`
	OccurredAt   int64  `json:"occurred_at"`
}
