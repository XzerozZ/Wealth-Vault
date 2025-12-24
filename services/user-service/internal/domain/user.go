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
