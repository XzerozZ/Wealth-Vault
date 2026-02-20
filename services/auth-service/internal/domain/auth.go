package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuthAccount struct {
	ID                uuid.UUID `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	UserID            uuid.UUID `json:"u_id"  gorm:"not null;index"`
	Provider          string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_provider_account;uniqueIndex:idx_email_provider"` // 'local', 'google'
	ProviderAccountID string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_provider_account"`
	Email             string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_email_provider"`
	Password          string    `gorm:"type:varchar(255)"`
	IsEmailVerified   bool      `json:"is_email_verified" gorm:"default:false"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type AuthOTP struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"u_id"  gorm:"not null"`
	OTP       string    `gorm:"type:varchar(6);not null"`
	ExpiredAt time.Time `gorm:"not null"`
}

type AuthSession struct {
	ID               uuid.UUID `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	UserID           uuid.UUID `json:"u_id"  gorm:"not null;index"`
	AccessToken      string    `gorm:"type:text;not null;unique"`
	RefreshToken     string    `gorm:"type:text;not null;uniqueIndex"`
	ExpiresAt        time.Time `gorm:"not null"`
	RefreshExpiresAt time.Time `gorm:"not null"`
	Revoked          bool      `gorm:"default:false"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
