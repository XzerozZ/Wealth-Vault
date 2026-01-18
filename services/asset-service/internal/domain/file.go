package domain

import (
	"time"

	"github.com/google/uuid"
)

type FileAssociate struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EntityID   uuid.UUID `json:"entity_id" gorm:"type:uuid;not null;index:idx_entity_type_id"`
	EntityType string    `json:"entity_type" gorm:"not null;index:idx_entity_type_id"`
	Link       string    `json:"link"  gorm:"not null"`
	FileType   string    `json:"file_type" gorm:"not null"`
	UserID     uuid.UUID `json:"u_id"  gorm:"type:uuid;not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
