package domain

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID           uuid.UUID     `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	GroupName    string        `json:"group_name" gorm:"not null"`
	GroupProfile string        `json:"group_profile" gorm:"not null"`
	CreatedBy    uuid.UUID     `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Creator      User          `gorm:"foreignKey:CreatedBy"`
	Members      []GroupMember `gorm:"foreignKey:GroupID"`
}

type GroupMember struct {
	GroupID  uuid.UUID `json:"group_id" gorm:"primaryKey"`
	UserID   uuid.UUID `json:"user_id" gorm:"primaryKey;type:uuid"`
	Role     string    `json:"role" gorm:"default:'member'"` // เช่น 'admin', 'member'
	JoinedAt time.Time `json:"joined_at" gorm:"autoCreateTime"`
	Group    Group     `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User     User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type GroupLog struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	GroupID   string    `json:"group_id" gorm:"not null"`
	Messages  string    `json:"messages" gorm:"not null"`
	CreatedBy string    `json:"created_by" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Group Group `gorm:"foreignKey:GroupID"`
	User  User  `gorm:"foreignKey:CreatedBy"`
}

type GroupItem struct {
	ID         uuid.UUID `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	GroupID    string    `json:"group_id" gorm:"not null;uniqueIndex:idx_group_item"`
	EntityType string    `json:"entity_type" gorm:"not null;uniqueIndex:idx_group_item"`
	EntityID   string    `json:"entity_id" gorm:"not null;uniqueIndex:idx_group_item"`
	CreatedBy  string    `json:"created_by" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Group Group `gorm:"foreignKey:GroupID"`
	User  User  `gorm:"foreignKey:CreatedBy"`
}
