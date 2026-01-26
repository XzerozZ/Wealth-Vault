package domain

import (
	"time"

	"github.com/google/uuid"
)

type Location struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Address     string    `gorm:"not null"`
	Subdistrict string    `gorm:"not null"`
	District    string    `gorm:"not null"`
	Province    string    `gorm:"not null"`
	PostalCode  string    `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
