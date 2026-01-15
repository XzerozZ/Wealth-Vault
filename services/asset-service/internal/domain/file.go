package domain

import "time"

type FileAssociate struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	EntityID   string    `json:"entity_id" gorm:"not null;index:idx_entity_type_id"`
	EntityType string    `json:"entity_type" gorm:"not null;index:idx_entity_type_id"`
	Link       string    `json:"link"  gorm:"not null"`
	FileType   string    `json:"file_type" gorm:"not null"`
	UserID     string    `json:"u_id"  gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
