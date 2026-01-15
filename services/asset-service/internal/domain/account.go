package domain

import "time"

type Account struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	Bank          string    `json:"bank" gorm:"not null"`
	Accountnumber string    `json:"acc_num" gorm:"not null"`
	Accounttype   string    `json:"acc_type" gorm:"not null"`
	Amount        float64   `json:"amount"`
	UserID        string    `json:"u_id"  gorm:"not null"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
