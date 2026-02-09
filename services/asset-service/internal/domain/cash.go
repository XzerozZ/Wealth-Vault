package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cash struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"not null"`
	Amount      float64   `gorm:"not null"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt  `gorm:"index"`
	Files       []FileAssociate `gorm:"polymorphic:Entity;polymorphicValue:cash"`
}
