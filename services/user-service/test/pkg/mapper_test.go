package utils_test

import (
	"testing"
	"time"

	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/pkg/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProtoMappers(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	t.Run("ToUserProto - Success with all fields", func(t *testing.T) {
		birthday := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
		uid := uuid.New()

		d := &domain.User{
			ID:                 uid,
			Email:              "test@wealthvault.com",
			Firstname:          "John",
			Lastname:           "Doe",
			Username:           "johndoe",
			Profile:            "https://img.com/p.png",
			Phonenumber:        "0812345678",
			Birthday:           &birthday,
			AutoShareAge:       25,
			IsAutoShareEnabled: true,
			CreatedAt:          now,
			UpdatedAt:          now,
		}

		result := utils.ToUserProto(d)

		assert.NotNil(t, result)
		assert.Equal(t, uid.String(), result.Id)
		assert.Equal(t, "test@wealthvault.com", result.Email)
		assert.Equal(t, int32(25), result.Sharedage)
		assert.True(t, result.Sharedenabled)
		assert.Equal(t, birthday.Unix(), result.Birthday.AsTime().Unix())
		assert.Equal(t, now.Unix(), result.CreatedAt.AsTime().Unix())
	})

	t.Run("ToUserProto - Success with Nil Birthday", func(t *testing.T) {
		d := &domain.User{
			ID:       uuid.New(),
			Birthday: nil,
		}

		result := utils.ToUserProto(d)

		assert.NotNil(t, result)
		assert.Nil(t, result.Birthday)
	})

	t.Run("ToGroupProto - Success", func(t *testing.T) {
		gid := uuid.New()
		uid := uuid.New()

		g := &domain.Group{
			ID:           gid,
			GroupName:    "Family Vault",
			GroupProfile: "https://img.com/g.png",
			CreatedBy:    uid,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		result := utils.ToGroupProto(g)

		assert.NotNil(t, result)
		assert.Equal(t, gid.String(), result.Id)
		assert.Equal(t, "Family Vault", result.Name)
		assert.Equal(t, uid.String(), result.UserId)
		assert.Equal(t, int64(0), result.MemberCount)
	})

	t.Run("Input Nil - Should return Nil", func(t *testing.T) {
		assert.Nil(t, utils.ToUserProto(nil))
		assert.Nil(t, utils.ToGroupProto(nil))
	})
}
