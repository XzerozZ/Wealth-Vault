package domain

import (
	"time"

	"github.com/google/uuid"
)

type AssetIDResult struct {
	Type string
	ID   string
}

type AssetBasicInfo struct {
	ID   uuid.UUID
	Name string
}

type AssetSummary struct {
	ID        uuid.UUID
	Type      string
	Name      string
	Value     float64
	Files     []FileAssociate
	CreatedAt time.Time
}

type NetWorthOverview struct {
	TotalAssets      float64
	TotalLiabilities float64
}
