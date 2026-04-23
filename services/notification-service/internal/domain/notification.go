package domain

import (
	"time"

	"github.com/google/uuid"
)

type WSMessage struct {
	Type    string      `json:"type"`
	Event   string      `json:"event,omitempty"`
	Payload interface{} `json:"payload"`
}

type Notification struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EntityType string    `gorm:"not null"`
	EntityID   uuid.UUID `gorm:"not null"`
	Receiver   uuid.UUID `gorm:"not null"`
	SenderID   *uuid.UUID
	Channel    string `gorm:"type:varchar(20);not null"`
	Message    string `gorm:"not null"`
	Metadata   string `gorm:"type:jsonb;default:'{}'"`
	CreatedAt  time.Time
	IsRead     bool ` gorm:"default:false"`
}
