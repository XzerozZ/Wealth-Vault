package domain

import "time"

type Cash struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Value       float64   `json:"value" gorm:"not null"`
	Description string    `json:"desc"`
	UserID      string    `json:"u_id"  gorm:"not null;index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
