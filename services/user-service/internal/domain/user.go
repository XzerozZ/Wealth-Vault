package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email       string    `gorm:"not null;uniqueIndex"`
	Firstname   string
	Lastname    string
	Username    string `gorm:"not null"`
	Profile     string
	Phonenumber string
	Birthday    *time.Time `gorm:"type:date"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FriendList struct {
	UserID   string `json:"user" gorm:"primaryKey;not null"`
	FriendID string `json:"friend" gorm:"primaryKey;not null"`
	User     User   `gorm:"foreignKey:UserID"`
	Friend   User   `gorm:"foreignKey:FriendID"`
}
