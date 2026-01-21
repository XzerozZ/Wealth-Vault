package helper

import (
	"context"
	"fmt"
	"wealth-vault/asset-service/internal/domain"
	repo "wealth-vault/asset-service/internal/repository/interface"

	"github.com/google/uuid"
)

type StorageDeleter interface {
	Delete(url string) error
}

func DeleteFilesAsync(storage StorageDeleter, fileURLs []string) {
	if len(fileURLs) == 0 {
		return
	}

	go func(urls []string) {
		for _, url := range urls {
			if err := storage.Delete(url); err != nil {
				fmt.Printf("⚠️ [AsyncDelete] Failed to delete file %s: %v\n", url, err)
			}
		}
	}(fileURLs)
}

func SyncEntityFiles(ctx context.Context, repo repo.FileRepository, storage StorageDeleter, params domain.FileSyncParams) error {
	if len(params.DeleteFileIDs) > 0 {
		var fileUUIDs []uuid.UUID
		for _, idStr := range params.DeleteFileIDs {
			if parsedID, err := uuid.Parse(idStr); err == nil {
				fileUUIDs = append(fileUUIDs, parsedID)
			}
		}

		if len(fileUUIDs) > 0 {
			filesToDelete, err := repo.GetFilesByIDs(params.Ctx, fileUUIDs)
			if err == nil {
				var validLinks []string
				var validIDs []string

				for _, f := range filesToDelete {
					if f.UserID == params.UserID {
						validLinks = append(validLinks, f.Link)
						validIDs = append(validIDs, f.ID.String())
					}
				}

				if len(validIDs) > 0 {
					go DeleteFilesAsync(storage, validLinks)
					err = repo.DeleteFiles(params.Ctx, validIDs, params.EntityID, params.UserID)
					if err != nil {
						return fmt.Errorf("failed to delete files: %w", err)
					}
				}
			}
		}
	}

	if len(params.NewFiles) > 0 {
		var filesToCreate []domain.FileAssociate

		for _, f := range params.NewFiles {
			if f.Url == "" {
				continue
			}

			filesToCreate = append(filesToCreate, domain.FileAssociate{
				ID:         uuid.New(),
				EntityID:   params.EntityID,
				EntityType: params.EntityType,
				UserID:     params.UserID,
				Link:       f.Url,
				FileType:   f.FileType,
			})
		}

		if len(filesToCreate) > 0 {
			if err := repo.CreateFiles(params.Ctx, filesToCreate); err != nil {
				return fmt.Errorf("failed to create files: %w", err)
			}
		}
	}

	return nil
}
