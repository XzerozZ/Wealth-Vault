package domain

import (
	"time"
)

type Notification struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	EntityType string    `json:"entity_type" gorm:"not null"`
	EntityID   string    `json:"entity_id" gorm:"not null"`
	Receiver   string    `json:"receiver" gorm:"not null"`
	SenderID   string    `json:"sender_id"`
	Channel    string    `json:"channel" gorm:"type:varchar(20);not null"`
	CreatedAt  time.Time `json:"created_at"`

	Sender User `gorm:"foreignKey:SenderID"`
}
