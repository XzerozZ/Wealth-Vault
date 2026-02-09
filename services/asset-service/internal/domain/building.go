package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BuildingType string

const (
	BuildingTypeCondo      BuildingType = "Condo"
	BuildingTypeHouse      BuildingType = "House"
	BuildingTypeTownHome   BuildingType = "Townhome"
	BuildingTypeCommercial BuildingType = "Commercial"
)

type Building struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID    `gorm:"type:uuid;not null;index"`
	Type        BuildingType `gorm:"type:varchar(50);not null;index"`
	Name        string       `gorm:"not null"`
	Area        float64      `gorm:"not null;default:0"`
	Amount      float64      `gorm:"not null;default:0"`
	Description string
	LocationID  uuid.UUID   `gorm:"type:uuid;not null"`
	Location    Location    `gorm:"foreignKey:LocationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Lands       []Land      `gorm:"many2many:building_land;joinForeignKey:house_id;joinReferences:land_id"`
	Insurances  []Insurance `gorm:"many2many:building_insurance;joinForeignKey:house_id;joinReferences:ins_id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt  `gorm:"index"`
	Files       []FileAssociate `gorm:"polymorphic:Entity;polymorphicValue:building"`
}
