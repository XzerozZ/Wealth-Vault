package domain

import "time"

type Group struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	GroupName string    `json:"group_name" gorm:"not null"`
	CreatedBy string    `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Creator User `gorm:"foreignKey:CreatedBy"`
}

type GroupLog struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	GroupID   string    `json:"group_id" gorm:"not null"`
	Messages  string    `json:"messages" gorm:"not null"`
	CreatedBy string    `json:"created_by" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Group Group `gorm:"foreignKey:GroupID"`
	User  User  `gorm:"foreignKey:CreatedBy"`
}

type GroupItem struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	GroupID    string    `json:"group_id" gorm:"not null;uniqueIndex:idx_group_item"`
	EntityType string    `json:"entity_type" gorm:"not null;uniqueIndex:idx_group_item"`
	EntityID   string    `json:"entity_id" gorm:"not null;uniqueIndex:idx_group_item"`
	CreatedBy  string    `json:"created_by" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Group Group `gorm:"foreignKey:GroupID"`
	User  User  `gorm:"foreignKey:CreatedBy"`
}
