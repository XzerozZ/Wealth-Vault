package domain

import "time"

type Cash struct {
	ID          string          `json:"id" gorm:"primaryKey"`
	Name        string          `json:"name"`
	Value       float64         `json:"value" gorm:"not null"`
	Description string          `json:"desc"`
	UserID      string          `json:"u_id"  gorm:"not null"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Files       []FileAssociate `gorm:"polymorphic:Entity;polymorphicValue:cash"`
}

type UpdateCashInput struct {
	ID            string
	Name          string
	Value         float64
	Description   string
	UserID        string
	UpdateMask    []string
	NewFiles      []FileAssociate
	DeleteFileIDs []string
}
