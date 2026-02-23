package helper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"wealth-vault/user-service/internal/domain"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils/helper"
	mock_storage "wealth-vault/user-service/test/mock/storage"
)

func TestApplyUpdateUserFields(t *testing.T) {
	t.Run("Update profile and delete old file", func(t *testing.T) {
		mockStorage := new(mock_storage.MockSupabaseStorage)
		oldProfile := "https://supabase.com/old_user.png"
		user := &domain.User{
			Username: "old_nick",
			Profile:  oldProfile,
		}

		req := &pb.UpdateUserRequest{
			Username:   "new_nick",
			Profile:    "https://supabase.com/new_user.png",
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"profile"}},
		}

		mockStorage.On("Delete", oldProfile).Return(nil).Once()

		mask, err := helper.ApplyUpdateUserFields(req, mockStorage, user)

		assert.NoError(t, err)
		assert.Equal(t, "old_nick", user.Username)
		assert.Equal(t, req.Profile, user.Profile)
		assert.Contains(t, mask, "Profile")

		assert.Eventually(t, func() bool {
			return len(mockStorage.Calls) == 1
		}, 1*time.Second, 50*time.Millisecond, "Old profile should be deleted in background")
	})

	t.Run("Update Birthday and Boolean fields", func(t *testing.T) {
		mockStorage := new(mock_storage.MockSupabaseStorage)
		user := &domain.User{}

		now := time.Now()
		sharedEnabled := true
		var sharedAge int32 = 25

		req := &pb.UpdateUserRequest{
			Birthday:      timestamppb.New(now),
			Sharedenabled: &sharedEnabled,
			Sharedage:     &sharedAge,
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{
				"birthday",
				"shared_age",
				"share_enabled",
			}},
		}

		mask, err := helper.ApplyUpdateUserFields(req, mockStorage, user)

		assert.NoError(t, err)

		assert.True(t, user.IsAutoShareEnabled, "IsAutoShareEnabled should be true")
		assert.Equal(t, 25, user.AutoShareAge, "AutoShareAge should be 25")

		assert.Contains(t, mask, "IsAutoShareEnabled")
		assert.Contains(t, mask, "AutoShareAge")
		assert.Contains(t, mask, "Birthday")
	})

	t.Run("No update when mask paths don't match", func(t *testing.T) {
		mockStorage := new(mock_storage.MockSupabaseStorage)
		user := &domain.User{Firstname: "John"}

		req := &pb.UpdateUserRequest{
			Firstname:  "Jane",
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"other_field"}},
		}

		mask, err := helper.ApplyUpdateUserFields(req, mockStorage, user)

		assert.NoError(t, err)
		assert.Equal(t, "John", user.Firstname)
		assert.Empty(t, mask)
	})
}
