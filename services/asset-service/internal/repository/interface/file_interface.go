package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
)

type FileRepository interface {
	CreateFiles(ctx context.Context, files []domain.FileAssociate) error
	DeleteFiles(ctx context.Context, fileIDs []string, entityType uuid.UUID, userID uuid.UUID) error
	GetFilesByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.FileAssociate, error)
}
