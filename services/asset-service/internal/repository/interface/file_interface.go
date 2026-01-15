package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"
)

type FileRepository interface {
	DeleteFiles(ctx context.Context, fileIDs []string, entityType string, userID string) error
	CreateFiles(ctx context.Context, files []domain.FileAssociate) error
}
