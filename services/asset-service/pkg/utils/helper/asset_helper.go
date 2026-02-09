package helper

import (
	"context"
	"log"
	"time"
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/user"

	"github.com/google/uuid"
)

func CleanupAssetResource(
	ctx context.Context,
	entityID uuid.UUID,
	files []domain.FileAssociate,
	storage StorageDeleter,
	userClient pb.UserServiceClient,
	hardDeleteFunc func(uuid.UUID) error,
) {

	if len(files) > 0 {
		var urls []string
		for _, f := range files {
			urls = append(urls, f.Link)
		}

		DeleteFilesAsync(storage, urls)
	}

	if err := hardDeleteFunc(entityID); err != nil {
		log.Printf("❌ Failed to hard delete asset %s: %v", entityID, err)
		return
	}

	notifyCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := userClient.DeleteAllReferencesByEntityID(notifyCtx, &pb.DeleteByEntityRequest{
		EntityId: entityID.String(),
	})

	if err != nil {
		log.Printf("⚠️ Warning: Asset %s deleted, but failed to notify Share Service: %v", entityID, err)
	} else {
		log.Printf("🧹 Cleanup Complete: Asset %s and all shared references deleted.", entityID)
	}
}
