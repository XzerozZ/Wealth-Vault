package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AssetType string

const (
	AssetTypeCash       AssetType = "CASH"
	AssetTypeBank       AssetType = "BANK"
	AssetTypeInvestment AssetType = "INVESTMENT"
	AssetTypeRealEstate AssetType = "REAL_ESTATE"
	AssetTypeInsurance  AssetType = "INSURANCE"
)

type Asset struct {
	ID                   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID               uuid.UUID `gorm:"type:uuid;not null;index"`
	Type                 AssetType `gorm:"type:varchar(50);not null;index"`
	Name                 string    `gorm:"not null"`
	Amount               float64   `gorm:"not null;default:0"`
	IsIncludedInNetWorth *bool     `gorm:"default:true;not null"`
	Description          string
	Details              datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Files                []FileAssociate `gorm:"polymorphic:Entity;polymorphicValue:asset"`
}
