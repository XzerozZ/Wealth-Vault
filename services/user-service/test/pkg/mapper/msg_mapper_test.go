package mapper_test

import (
	"testing"
	"time"

	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/pkg/utils/mapper"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMessageMapper(t *testing.T) {
	myID := uuid.New()
	otherID := uuid.New()
	msgID := uuid.New()
	now := time.Now()

	t.Run("ToGroupMessageProto - Success (IsMe = true)", func(t *testing.T) {
		d := domain.GroupMessage{
			ID:        msgID,
			SenderID:  myID,
			MsgType:   "TEXT",
			Content:   "Hello Group",
			CreatedAt: now,
			Sender: &domain.User{
				Username: "Nayme",
				Profile:  "https://image.com/me.png",
			},
		}

		result := mapper.ToGroupMessageProto(d, myID.String())

		assert.NotNil(t, result)
		assert.Equal(t, msgID.String(), result.Id)
		assert.Equal(t, "Nayme", result.SenderName)
		assert.Equal(t, "https://image.com/me.png", result.SenderImage)
		assert.True(t, result.IsMe)
		assert.Equal(t, now.Unix(), result.CreatedAt.AsTime().Unix())
	})

	t.Run("ToGroupMessageProto - Success (IsMe = false, Sender is Nil)", func(t *testing.T) {
		d := domain.GroupMessage{
			ID:        msgID,
			SenderID:  otherID,
			MsgType:   "IMAGE",
			Content:   "img_url",
			CreatedAt: now,
			Sender:    nil,
		}

		result := mapper.ToGroupMessageProto(d, myID.String())

		assert.NotNil(t, result)
		assert.Equal(t, "Unknown", result.SenderName)
		assert.Equal(t, "", result.SenderImage)
		assert.False(t, result.IsMe)
	})

	t.Run("ToPrivateMessageProto - Success", func(t *testing.T) {
		d := domain.PrivateMessage{
			ID:        msgID,
			SenderID:  otherID,
			MsgType:   "TEXT",
			Content:   "Hi there",
			CreatedAt: now,
			Sender: &domain.User{
				Username: "Friend",
				Profile:  "https://image.com/friend.png",
			},
		}

		result := mapper.ToPrivateMessageProto(d, myID.String())

		assert.NotNil(t, result)
		assert.Equal(t, "Friend", result.SenderName)
		assert.False(t, result.IsMe)
		assert.Equal(t, "Hi there", result.Content)
	})
}
