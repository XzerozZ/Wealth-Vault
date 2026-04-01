package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Land struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"not null"`
	DeedNum     string    `gorm:"not null;unique"`
	Area        float64   `gorm:"not null;default:0"`
	Amount      float64   `gorm:"not null;default:0"`
	Description string
	LocationID  uuid.UUID  `gorm:"type:uuid;not null"`
	Location    Location   `gorm:"foreignKey:LocationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Buildings   []Building `gorm:"many2many:building_land"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt  `gorm:"index"`
	Files       []FileAssociate `gorm:"polymorphic:Entity;polymorphicValue:land"`
}
