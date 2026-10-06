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

// LIFF identities are separate from legacy provider accounts so existing Kotlin
// LINE links remain untouched during the migration.
type LineIdentity struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ChannelID string    `gorm:"type:varchar(80);not null;uniqueIndex:idx_line_channel_subject"`
	Subject   string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_line_channel_subject"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt time.Time
}

type LineChallenge struct {
	Digest     string    `gorm:"type:varchar(64);primaryKey"`
	ChannelID  string    `gorm:"type:varchar(80);not null"`
	Subject    string    `gorm:"type:varchar(255);not null"`
	ExpiresAt  time.Time `gorm:"not null;index"`
	ConsumedAt *time.Time
	Attempts   int `gorm:"not null;default:0"`
	CreatedAt  time.Time
}

type LineAuthResult struct {
	Status               string
	ChallengeID          string
	UserID               string
	DisplayName          string
	AccessToken          string
	RefreshToken         string
	AccessTokenExpiresIn int32
}

type LineChallengeInput struct {
	ChallengeID string
	ChannelID   string
	Action      string
	Email       string
	Password    string
	Username    string
	FirstName   string
	LastName    string
}
