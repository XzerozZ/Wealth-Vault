package usecase_test

import (
	"context"
	"testing"
	"time"

	"wealth-vault/notification-service/internal/domain"
	"wealth-vault/notification-service/internal/usecase"
	mock_repo "wealth-vault/notification-service/test/mock/repository"
	mock_socket "wealth-vault/notification-service/test/mock/socket"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNotificationUsecase(t *testing.T) {
	t.Run("GetHistory - Success", func(t *testing.T) {
		repo := new(mock_repo.MockNotificationRepository)
		hub := new(mock_socket.MockSocketHub)
		uc := usecase.NewNotificationUsecase(repo, hub)

		ctx := context.Background()
		uid := uuid.New()

		expected := []domain.Notification{
			{ID: uuid.New(), Message: "Test"},
		}

		repo.On("GetByReceiver", ctx, uid).Return(expected, nil).Once()

		res, err := uc.GetHistory(ctx, uid)

		assert.NoError(t, err)
		assert.Equal(t, expected, res)

		repo.AssertExpectations(t)
	})

	t.Run("HandleFriendRequest - Success", func(t *testing.T) {
		repo := new(mock_repo.MockNotificationRepository)
		hub := new(mock_socket.MockSocketHub)
		uc := usecase.NewNotificationUsecase(repo, hub)

		ctx := context.Background()

		targetID := uuid.New()
		senderID := uuid.New()

		evt := domain.FriendRequestEvent{
			TargetID:   targetID.String(),
			SenderID:   senderID.String(),
			SenderName: "Alice",
			OccurredAt: time.Now().Unix(),
		}

		repo.
			On("CreateNotification", ctx, mock.MatchedBy(func(n *domain.Notification) bool {
				return n.Receiver == targetID &&
					n.EntityType == "FRIEND_REQUEST"
			})).
			Return(nil).
			Once()

		hub.
			On("Emit", targetID.String(), mock.MatchedBy(func(msg domain.WSMessage) bool {
				return msg.Type == "NOTIFICATION" &&
					msg.Payload.(*domain.Notification).EntityType == "FRIEND_REQUEST"
			})).
			Return().
			Once()

		err := uc.HandleFriendRequest(ctx, evt)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
		hub.AssertExpectations(t)
	})

	t.Run("HandleGroupMemberAdded - Multi-User and Skip Sender", func(t *testing.T) {
		repo := new(mock_repo.MockNotificationRepository)
		hub := new(mock_socket.MockSocketHub)
		uc := usecase.NewNotificationUsecase(repo, hub)

		ctx := context.Background()

		groupID := uuid.New()
		targetID := uuid.New()
		senderID := uuid.New()

		evt := domain.GroupMemberAddedEvent{
			GroupID:       groupID.String(),
			TargetUserIDs: []string{targetID.String(), senderID.String()},
			SenderID:      senderID.String(),
		}

		repo.
			On("CreateNotification", ctx, mock.MatchedBy(func(n *domain.Notification) bool {
				return n.Receiver == targetID
			})).
			Return(nil).
			Once()

		hub.
			On("Emit", targetID.String(), mock.Anything).
			Return().
			Once()

		hub.
			On("BroadcastToGroup", groupID.String(), mock.MatchedBy(func(msg domain.WSMessage) bool {
				return msg.Event == "MEMBER_ADDED"
			})).
			Return().
			Once()

		err := uc.HandleGroupMemberAdded(ctx, evt)

		assert.NoError(t, err)

		repo.AssertExpectations(t)
		hub.AssertExpectations(t)
	})

	t.Run("HandleMemberRemoved - Success", func(t *testing.T) {
		repo := new(mock_repo.MockNotificationRepository)
		hub := new(mock_socket.MockSocketHub)
		uc := usecase.NewNotificationUsecase(repo, hub)

		ctx := context.Background()

		targetID := uuid.New()
		groupID := uuid.New()

		evt := domain.MemberRemovedEvent{
			TargetID:  targetID.String(),
			GroupID:   groupID.String(),
			GroupName: "Family",
		}

		repo.
			On("CreateNotification", ctx, mock.Anything).
			Return(nil).
			Once()

		hub.
			On("Emit", targetID.String(), mock.MatchedBy(func(msg domain.WSMessage) bool {
				return msg.Type == "NOTIFICATION"
			})).
			Return().
			Once()

		hub.
			On("Emit", targetID.String(), mock.MatchedBy(func(msg domain.WSMessage) bool {
				return msg.Type == "DATA_UPDATE" &&
					msg.Event == "YOU_ARE_REMOVED"
			})).
			Return().
			Once()

		err := uc.HandleMemberRemoved(ctx, evt)

		assert.NoError(t, err)

		repo.AssertExpectations(t)
		hub.AssertExpectations(t)
	})

	t.Run("HandleInsuranceExpiring - Invalid UUID", func(t *testing.T) {
		repo := new(mock_repo.MockNotificationRepository)
		hub := new(mock_socket.MockSocketHub)
		uc := usecase.NewNotificationUsecase(repo, hub)

		ctx := context.Background()

		evt := domain.InsuranceExpiringEvent{
			UserID: "not-a-uuid",
		}

		err := uc.HandleInsuranceExpiring(ctx, evt)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user id")

		repo.AssertNotCalled(t, "CreateNotification", mock.Anything, mock.Anything)
	})

	t.Run("HandleAccessGranted - Check Payload", func(t *testing.T) {
		repo := new(mock_repo.MockNotificationRepository)
		hub := new(mock_socket.MockSocketHub)
		uc := usecase.NewNotificationUsecase(repo, hub)

		ctx := context.Background()

		targetID := uuid.New()
		groupID := uuid.New()

		evt := domain.AccessGrantedEvent{
			TargetUserID: targetID.String(),
			GroupID:      groupID.String(),
			ItemCount:    5,
			GrantorName:  "Admin",
		}

		repo.
			On("CreateNotification", ctx, mock.Anything).
			Return(nil).
			Once()

		hub.
			On("Emit", targetID.String(), mock.MatchedBy(func(msg domain.WSMessage) bool {
				return msg.Type == "NOTIFICATION"
			})).
			Return().
			Once()

		hub.
			On("Emit", targetID.String(), mock.MatchedBy(func(msg domain.WSMessage) bool {
				if msg.Type != "DATA_UPDATE" || msg.Event != "ACCESS_GRANTED" {
					return false
				}

				p := msg.Payload.(map[string]interface{})
				return p["group_id"] == groupID.String() &&
					p["count"] == 5
			})).
			Return().
			Once()

		err := uc.HandleAccessGranted(ctx, evt)

		assert.NoError(t, err)

		repo.AssertExpectations(t)
		hub.AssertExpectations(t)
	})
}
