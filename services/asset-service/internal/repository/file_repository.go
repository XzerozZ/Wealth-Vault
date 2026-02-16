package repository

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) DeleteFiles(ctx context.Context, fileIDs []uuid.UUID) error {
	if err := r.db.WithContext(ctx).Unscoped().
		Where("id IN ?", fileIDs).
		Delete(&domain.FileAssociate{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *FileRepository) CreateFiles(ctx context.Context, files []domain.FileAssociate) error {
	var dbFiles []domain.FileAssociate
	for _, f := range files {
		dbFiles = append(dbFiles, domain.FileAssociate{
			EntityID:   f.EntityID,
			EntityType: f.EntityType,
			Link:       f.Link,
			FileType:   f.FileType,
			UserID:     f.UserID,
		})
	}

	if err := r.db.WithContext(ctx).Create(&dbFiles).Error; err != nil {
		return err
	}

	return nil
}

func (r *FileRepository) GetFilesByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.FileAssociate, error) {
	var files []domain.FileAssociate
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&files).Error; err != nil {
		return nil, err
	}

	return files, nil
}
