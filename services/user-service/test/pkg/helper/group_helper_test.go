package helper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"wealth-vault/user-service/internal/domain"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils/helper"
	mock_storage "wealth-vault/user-service/test/mock/storage"
)

func TestApplyUpdateGroupFields(t *testing.T) {
	t.Run("Update name only using field mask", func(t *testing.T) {
		mockStorage := new(mock_storage.MockSupabaseStorage)
		group := &domain.Group{GroupName: "Old Name", GroupProfile: "old_url"}

		req := &pb.UpdateGroupRequest{
			Name:       "New Name",
			Profile:    "new_url",
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"group_name"}},
		}

		mask, err := helper.ApplyUpdateGroupFields(req, mockStorage, group)

		assert.NoError(t, err)
		assert.Equal(t, "New Name", group.GroupName)
		assert.Equal(t, "old_url", group.GroupProfile)
		assert.Contains(t, mask, "GroupName")
		assert.NotContains(t, mask, "GroupProfile")
	})

	t.Run("Full update when mask is nil", func(t *testing.T) {
		mockStorage := new(mock_storage.MockSupabaseStorage)
		group := &domain.Group{GroupName: "Old Name", GroupProfile: "old_url"}

		req := &pb.UpdateGroupRequest{
			Name:    "New Name",
			Profile: "new_url",
		}

		mockStorage.On("Delete", "old_url").Return(nil).Once()

		_, err := helper.ApplyUpdateGroupFields(req, mockStorage, group)

		assert.NoError(t, err)

		time.Sleep(50 * time.Millisecond)

		assert.Eventually(t, func() bool {
			return mockStorage.AssertExpectations(t)
		}, 1*time.Second, 50*time.Millisecond)
	})

	t.Run("Trigger file deletion when profile changes", func(t *testing.T) {
		mockStorage := new(mock_storage.MockSupabaseStorage)
		oldProfile := "https://supabase.com/old.png"
		group := &domain.Group{GroupProfile: oldProfile}

		req := &pb.UpdateGroupRequest{
			Profile:    "https://supabase.com/new.png",
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"profile"}},
		}

		mockStorage.On("Delete", oldProfile).Return(nil).Once()

		_, err := helper.ApplyUpdateGroupFields(req, mockStorage, group)
		assert.NoError(t, err)

		assert.Eventually(t, func() bool {
			callCount := 0
			for _, call := range mockStorage.Calls {
				if call.Method == "Delete" {
					callCount++
				}
			}
			return callCount == 1
		}, 1*time.Second, 50*time.Millisecond, "Storage Delete should be called once in background")

		mockStorage.AssertExpectations(t)
	})
}
