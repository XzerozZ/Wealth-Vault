package repository

import (
	"context"

	"github.com/google/uuid"
)

type AssetRepository interface {
	CheckExists(ctx context.Context, entityType string, id uuid.UUID, uid uuid.UUID) (bool, error)
}
