package domain

import (
	"time"

	"github.com/google/uuid"
)

type BankType string

const (
	BankTypeSavings      BankType = "Savings"
	BankTypeCurrent      BankType = "Current"
	BankTypeFixedDeposit BankType = "Fixed_Deposit"
)

type Account struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"not null"`
	BankName    string    `gorm:"not null"`
	BankAccount string    `gorm:"not null;unique"`
	Type        BankType  `gorm:"type:varchar(50);not null;index"`
	Amount      float64   `gorm:"not null;default:0"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Files       []FileAssociate `gorm:"polymorphic:Entity;polymorphicValue:account"`
}
