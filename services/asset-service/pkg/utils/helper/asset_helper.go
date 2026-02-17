package helper

import (
	"context"
	"fmt"
	"log"
	"time"

	"wealth-vault/asset-service/internal/domain"
	repo "wealth-vault/asset-service/internal/repository/interface"
	pb "wealth-vault/asset-service/pkg/pb/proto/user"

	"github.com/google/uuid"
)

type AssetHelper interface {
	SyncFiles(ctx context.Context, params domain.FileSyncParams) error
	CleanupResource(ctx context.Context, entityID uuid.UUID, files []domain.FileAssociate, hardDeleteFunc func(uuid.UUID) error)
}

type StorageDeleter interface {
	Delete(url string) error
}

type RealAssetHelper struct {
	fileRepo   repo.FileRepository
	storage    StorageDeleter
	userClient pb.UserServiceClient
}

func NewAssetHelper(fr repo.FileRepository, sd StorageDeleter, uc pb.UserServiceClient) AssetHelper {
	return &RealAssetHelper{
		fileRepo:   fr,
		storage:    sd,
		userClient: uc,
	}
}

func (h *RealAssetHelper) SyncFiles(ctx context.Context, params domain.FileSyncParams) error {
	return SyncEntityFiles(ctx, h.fileRepo, h.storage, params)
}

func (h *RealAssetHelper) CleanupResource(
	ctx context.Context,
	entityID uuid.UUID,
	files []domain.FileAssociate,
	hardDeleteFunc func(uuid.UUID) error,
) {
	CleanupAssetResource(ctx, entityID, files, h.storage, h.userClient, hardDeleteFunc)
}

func DeleteFilesAsync(storage StorageDeleter, fileURLs []string) {
	if len(fileURLs) == 0 {
		return
	}

	go func() {
		for _, url := range fileURLs {
			if err := storage.Delete(url); err != nil {
				log.Printf("[AsyncDelete] Failed: %s | Error: %v", url, err)
			}
		}
	}()
}

func SyncEntityFiles(ctx context.Context, r repo.FileRepository, storage StorageDeleter, params domain.FileSyncParams) error {
	if len(params.DeleteFileIDs) > 0 {
		fileUUIDs := parseUUIDs(params.DeleteFileIDs)

		if len(fileUUIDs) > 0 {
			files, err := r.GetFilesByIDs(ctx, fileUUIDs)
			if err == nil {
				var validLinks []string
				var validIDs []string

				for _, f := range files {
					if f.UserID == params.UserID {
						validLinks = append(validLinks, f.Link)
						validIDs = append(validIDs, f.ID.String())
					}
				}

				if len(validIDs) > 0 {
					DeleteFilesAsync(storage, validLinks)
					if err := r.DeleteFiles(ctx, fileUUIDs); err != nil {
						log.Printf("❌ [Cleanup] DB delete failed: %v", err)
						return fmt.Errorf("db deletion failed: %w", err)
					}
				}
			}
		}
	}

	if len(params.NewFiles) > 0 {
		filesToCreate := make([]domain.FileAssociate, 0, len(params.NewFiles))
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
			if err := r.CreateFiles(ctx, filesToCreate); err != nil {
				return fmt.Errorf("db creation failed: %w", err)
			}
		}
	}
	return nil
}

func CleanupAssetResource(
	ctx context.Context,
	entityID uuid.UUID,
	files []domain.FileAssociate,
	storage StorageDeleter,
	userClient pb.UserServiceClient,
	hardDeleteFunc func(uuid.UUID) error,
) {
	if len(files) > 0 {
		urls := make([]string, len(files))
		for i, f := range files {
			urls[i] = f.Link
		}
		DeleteFilesAsync(storage, urls)
	}

	if err := hardDeleteFunc(entityID); err != nil {
		log.Printf("[Cleanup] DB Delete Failed: %s | %v", entityID, err)
		return
	}

	go func(id string) {
		notifyCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		_, err := userClient.DeleteAllReferencesByEntityID(notifyCtx, &pb.DeleteByEntityRequest{
			EntityId: id,
		})
		if err != nil {
			log.Printf("[Cleanup] Notify Failed: %s | %v", id, err)
		} else {
			log.Printf("[Cleanup] Success: %s", id)
		}
	}(entityID.String())
}

func parseUUIDs(ids []string) []uuid.UUID {
	parsed := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if u, err := uuid.Parse(id); err == nil {
			parsed = append(parsed, u)
		}
	}
	return parsed
}
