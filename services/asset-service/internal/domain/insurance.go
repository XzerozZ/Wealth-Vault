package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InsuranceType string

const (
	InsuranceTypeLife     InsuranceType = "Life"
	InsuranceTypeHealth   InsuranceType = "Health"
	InsuranceTypeAccident InsuranceType = "Accident"
	InsuranceTypeProperty InsuranceType = "Property"
	InsuranceTypeVehicle  InsuranceType = "Vehicle"
)

type Insurance struct {
	ID             uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID         uuid.UUID     `gorm:"type:uuid;not null;index"`
	Name           string        `gorm:"not null"`
	PolicyNumber   string        `gorm:"not null"`
	Type           InsuranceType `gorm:"type:varchar(50);not null;index"`
	CompanyName    string        `gorm:"not null"`
	CoveragePeriod float64       `gorm:"not null;default:0"`
	CoverageAmount float64       `gorm:"not null;default:0"`
	ConDate        *time.Time    `gorm:"type:date"`
	ExpDate        *time.Time    `gorm:"type:date"`
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt  `gorm:"index"`
	Files          []FileAssociate `gorm:"polymorphic:Entity;polymorphicValue:insurance"`
	Buildings      []Building      `gorm:"many2many:building_insurance;joinForeignKey:ins_id;joinReferences:house_id"`
}
