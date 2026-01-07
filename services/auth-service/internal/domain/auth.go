package domain

import "time"

type AuthAccount struct {
	ID                string `json:"id" gorm:"primaryKey"`
	UserID            string `json:"u_id"  gorm:"not null;index"`
	Provider          string `gorm:"type:varchar(50);not null;uniqueIndex:idx_provider_account"` // 'local', 'google'
	ProviderAccountID string `gorm:"type:varchar(255);uniqueIndex:idx_provider_account"`
	Email             string `gorm:"type:varchar(255)"`
	Password          string `gorm:"type:varchar(255)"`
	IsEmailVerified   bool   `json:"is_email_verified" gorm:"default:false"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type AuthOTP struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"u_id"  gorm:"not null"`
	OTP       string    `gorm:"type:varchar(6);not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

type AuthSession struct {
	ID               string    `json:"id" gorm:"primaryKey"`
	UserID           string    `json:"u_id"  gorm:"not null;index"`
	AccessToken      string    `gorm:"type:text;not null;unique"`
	RefreshToken     string    `gorm:"type:text;not null;unique"`
	ExpiresAt        time.Time `gorm:"not null"`
	RefreshExpiresAt time.Time `gorm:"not null"`
	Revoked          bool      `gorm:"default:false"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
