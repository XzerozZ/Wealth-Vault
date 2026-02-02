package domain

import (
	"time"

	"github.com/google/uuid"
)

type FriendItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OwnerID    uuid.UUID `gorm:"type:uuid;not null"`
	FriendID   uuid.UUID `gorm:"type:uuid;not null"`
	EntityType string
	EntityID   uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ShareAt    time.Time `gorm:"index"`
	User       User      `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Friend     User      `gorm:"foreignKey:FriendID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type EmailItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OwnerID    uuid.UUID `gorm:"type:uuid;not null"`
	Email      string    `gorm:"not null"`
	EntityType string
	EntityID   uuid.UUID
	CreatedAt  time.Time
	ShareAt    time.Time `gorm:"index"`
	IsSent     bool      `gorm:"default:false"`
	User       User      `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type GroupItem struct {
	ID         uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	GroupID    uuid.UUID `gorm:"not null;uniqueIndex:idx_group_item"`
	EntityType string    `gorm:"not null;uniqueIndex:idx_group_item"`
	EntityID   uuid.UUID `gorm:"not null;uniqueIndex:idx_group_item"`
	OwnerID    uuid.UUID `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ShareAt    time.Time         `gorm:"index"`
	Group      Group             `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User       User              `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Viewers    []GroupItemViewer `gorm:"foreignKey:GroupItemID;constraint:OnDelete:CASCADE;"`
}

type GroupItemViewer struct {
	GroupItemID uuid.UUID `gorm:"primaryKey;type:uuid"`
	ViewerID    uuid.UUID `gorm:"primaryKey;type:uuid"`
	GrantedAt   time.Time `gorm:"autoCreateTime"`
	GroupItem   GroupItem `gorm:"foreignKey:GroupItemID"`
	User        User      `gorm:"foreignKey:ViewerID"`
}

type ItemWithDetail struct {
	ID         string
	SharedBy   string
	SharedAt   time.Time
	Type       string
	Building   *BuildingPreview
	Land       *LandPreview
	Insurance  *InsurancePreview
	Investment *InvestmentPreview
	Account    *AccountPreview
	Cash       *CashPreview
	Liability  *LiabilityPreview
}
