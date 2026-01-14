package domain

import (
	"time"
)

type User struct {
	ID          string    `json:"user_id" gorm:"primaryKey"`
	Email       string    `json:"email" gorm:"not null;uniqueIndex"`
	Firstname   string    `json:"first_name"`
	Lastname    string    `json:"last_name"`
	Username    string    `json:"username" gorm:"not null"`
	Profile     string    `json:"profile"`
	Phonenumber string    `json:"phone_number"`
	Birthday    time.Time `json:"birthday" gorm:"type:date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FriendList struct {
	UserID   string `json:"user" gorm:"primaryKey;not null"`
	FriendID string `json:"friend" gorm:"primaryKey;not null"`

	User   User `gorm:"foreignKey:UserID"`
	Friend User `gorm:"foreignKey:FriendID"`
}
