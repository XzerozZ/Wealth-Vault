package domain

import (
	"time"

	"github.com/google/uuid"
)

type AssetSummary struct {
	ID        uuid.UUID
	Type      string
	Name      string
	Value     float64
	CreatedAt time.Time
}

type NetWorthOverview struct {
	TotalAssets      float64
	TotalLiabilities float64
}
