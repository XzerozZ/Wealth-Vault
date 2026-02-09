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
	Creator      User          `gorm:"foreignKey:CreatedBy;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Members      []GroupMember `gorm:"foreignKey:GroupID"`
}

type GroupMember struct {
	GroupID  uuid.UUID `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	UserID   uuid.UUID `json:"user_id" gorm:"primaryKey;type:uuid"`
	Role     string    `json:"role" gorm:"default:'member'"`
	JoinedAt time.Time `json:"joined_at" gorm:"autoCreateTime"`
	Group    Group     `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User     User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type GroupLog struct {
	ID        uuid.UUID `gorm:"primaryKey;default:gen_random_uuid()"`
	GroupID   uuid.UUID `gorm:"not null;index"`
	LogType   string    `gorm:"not null;default:'SYSTEM'"`
	Metadata  string    `json:"metadata" gorm:"type:jsonb;default:'{}'"`
	Messages  string    `gorm:"not null"`
	CreatedBy uuid.UUID `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User ` gorm:"foreignKey:CreatedBy;references:ID"`
}

type GroupMessage struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	GroupID   uuid.UUID `gorm:"type:uuid;not null;index"`
	SenderID  uuid.UUID `gorm:"type:uuid;not null;index"`
	MsgType   string    `gorm:"type:varchar(20);default:'TEXT';not null"` // Type:"ASSET_CARD"
	Content   string    `gorm:"type:text"`
	Metadata  string    `json:"metadata" gorm:"type:jsonb;default:'{}'"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
	UpdatedAt time.Time `json:"updated_at"`
	Sender    *User     `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
}
