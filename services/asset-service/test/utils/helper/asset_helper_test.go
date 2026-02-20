package helper_test

import (
	"context"
	"testing"
	"time"
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	userpb "wealth-vault/asset-service/pkg/pb/proto/user"
	"wealth-vault/asset-service/pkg/utils/helper"
	mock_client "wealth-vault/asset-service/test/mock/client"
	mock_helper "wealth-vault/asset-service/test/mock/helper"
	mock_repo "wealth-vault/asset-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSyncEntityFiles_DeleteSuccess(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mock_repo.MockFileRepository)
	mockStorage := new(mock_helper.MockStorage)

	userID := uuid.New()
	entityID := uuid.New()
	fileID := uuid.New()

	t.Run("Success - Sync New and Delete Old Files", func(t *testing.T) {
		params := domain.FileSyncParams{
			UserID:        userID,
			EntityID:      entityID,
			EntityType:    "account",
			DeleteFileIDs: []string{fileID.String()},
			NewFiles:      []*pb.FileInfo{{Url: "new-file.jpg", FileType: "image/jpeg"}},
		}

		oldFiles := []domain.FileAssociate{{ID: fileID, UserID: userID, Link: "old-file.jpg"}}
		mockRepo.On("GetFilesByIDs", ctx, []uuid.UUID{fileID}).Return(oldFiles, nil)

		mockStorage.On("Delete", "old-file.jpg").Return(nil)

		mockRepo.On("DeleteFiles", ctx, []uuid.UUID{fileID}).Return(nil)

		mockRepo.On("CreateFiles", ctx, mock.MatchedBy(func(files []domain.FileAssociate) bool {
			return len(files) == 1 && files[0].Link == "new-file.jpg"
		})).Return(nil)

		err := helper.SyncEntityFiles(ctx, mockRepo, mockStorage, params)

		assert.NoError(t, err)
		time.Sleep(50 * time.Millisecond)
		mockRepo.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})
}

func TestCleanupAssetResource(t *testing.T) {
	mockStorage := new(mock_helper.MockStorage)
	mockUserClient := new(mock_client.MockUserClient)
	entityID := uuid.New()
	files := []domain.FileAssociate{{Link: "file-to-delete.jpg"}}

	mockStorage.On("Delete", "file-to-delete.jpg").Return(nil)

	mockUserClient.On("DeleteAllReferencesByEntityID", mock.Anything, &userpb.DeleteByEntityRequest{
		EntityId: entityID.String(),
	}).Return(&userpb.DeleteByEntityResponse{}, nil)

	hardDeleteCalled := false
	hardDeleteFunc := func(id uuid.UUID) error {
		hardDeleteCalled = true
		assert.Equal(t, entityID, id)
		return nil
	}

	helper.CleanupAssetResource(context.Background(), entityID, files, mockStorage, mockUserClient, hardDeleteFunc)

	assert.True(t, hardDeleteCalled)

	time.Sleep(100 * time.Millisecond)
	mockUserClient.AssertExpectations(t)
}
