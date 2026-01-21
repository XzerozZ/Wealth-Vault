package domain

import (
	"time"

	"github.com/google/uuid"
)

type LiabilityType string

const (
	LiabilityTypeLoan    LiabilityType = "LOAN"    // หนี้
	LiabilityTypeExpense LiabilityType = "EXPENSE" // ค่าใช้จ่ายต่อเนื่อง
)

type Liability struct {
	ID           uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID       uuid.UUID     `gorm:"type:uuid;not null;index"`
	Type         LiabilityType `gorm:"type:varchar(50);not null"`
	Creditor     string        `gorm:"not null"`
	Name         string        `gorm:"not null"`
	Principal    float64       `gorm:"not null"`
	InterestRate float64
	StartAt      *time.Time
	EndAt        *time.Time
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Files        []FileAssociate `gorm:"polymorphic:Entity;polymorphicValue:liability"`
}
