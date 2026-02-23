package helper_test

import (
	"testing"
	"time"

	"wealth-vault/user-service/pkg/utils/helper"
	mock_storage "wealth-vault/user-service/test/mock/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteFilesAsync(t *testing.T) {
	t.Run("should call storage delete for each URL in background", func(t *testing.T) {
		mockStorage := new(mock_storage.MockSupabaseStorage)
		fileURLs := []string{
			"https://test.com/1.png",
			"https://test.com/2.png",
		}

		mockStorage.On("Delete", mock.Anything).Return(nil).Times(2)

		helper.DeleteFilesAsync(mockStorage, fileURLs)

		assert.Eventually(t, func() bool {
			return len(mockStorage.Calls) == 2
		}, 1*time.Second, 100*time.Millisecond)

		mockStorage.AssertExpectations(t)
	})

	t.Run("should return immediately if fileURLs is empty", func(t *testing.T) {
		mockStorage := new(mock_storage.MockSupabaseStorage)

		helper.DeleteFilesAsync(mockStorage, []string{})

		time.Sleep(50 * time.Millisecond)
		mockStorage.AssertNotCalled(t, "Delete", mock.Anything)
	})
}
