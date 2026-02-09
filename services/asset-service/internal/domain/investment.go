package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvestmentType string

const (
	InvestTypeStockTH    InvestmentType = "Stock_TH"
	InvestTypeStockUS    InvestmentType = "Stock_US"
	InvestTypeMutualFund InvestmentType = "Mutual_Fund"
	InvestTypeBond       InvestmentType = "Bond"
	InvestTypeCrypto     InvestmentType = "Crypto"
	InvestTypeGold       InvestmentType = "Gold"
)

type Investment struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index"`
	Name         string         `gorm:"not null"`
	Symbol       string         `gorm:"not null"`
	Type         InvestmentType `gorm:"type:varchar(50);not null;index"`
	BrokerName   string         `gorm:"not null"`
	Quantity     float64        `gorm:"not null;default:0"`
	CostPerPrice float64        `gorm:"not null;default:0"`
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt  `gorm:"index"`
	Files        []FileAssociate `gorm:"polymorphic:Entity;polymorphicValue:investment"`
}
