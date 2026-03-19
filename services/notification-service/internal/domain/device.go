package domain

import (
	"time"

	"github.com/google/uuid"
)

type Platform string

const (
	PlatformIOS     Platform = "IOS"
	PlatformAndroid Platform = "Android"
)

type DeviceToken struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Token      string    `gorm:"type:text;not null;uniqueIndex"`
	Platform   Platform  `gorm:"type:varchar(10);not null"`
	DeviceName string    `gorm:"type:varchar(100)" `
	IsActive   bool      `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type RegisterDeviceRequest struct {
	Token      string   `json:"token"  form:"token"`
	Platform   Platform `json:"platform" form:"platform"`
	DeviceName string   `json:"device_name"  form:"device_name"`
}
