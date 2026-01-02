package domain

import "time"

type User struct {
	ID          string    `json:"u_id" gorm:"primaryKey"`
	Firstname   string    `json:"first_name"  gorm:"not null"`
	Lastname    string    `json:"last_name" gorm:"not null"`
	Username    string    `json:"username" gorm:"not null"`
	Profile     string    `json:"profile"`
	Phonenumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FriendList struct {
	UserID   string `json:"user" gorm:"primaryKey;not null"`
	FriendID string `json:"friend" gorm:"primaryKey;not null"`
	User     User   `gorm:"foreignKey:UserID"`
	Friend   User   `gorm:"foreignKey:FriendID"`
}
