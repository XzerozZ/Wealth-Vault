package usecase_test

import (
	"context"
	"testing"
	"time"

	"wealth-vault/notification-service/internal/domain"
	"wealth-vault/notification-service/internal/usecase"
	pb "wealth-vault/notification-service/pkg/pb/proto/auth"
	mock_client "wealth-vault/notification-service/test/mock/client"
	mock_line "wealth-vault/notification-service/test/mock/line"
	mock_dispatch "wealth-vault/notification-service/test/mock/push_provider"
	mock_repo "wealth-vault/notification-service/test/mock/repository"
	mock_socket "wealth-vault/notification-service/test/mock/socket"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNotificationUsecase(t *testing.T) {
	t.Run("GetHistory - Success", func(t *testing.T) {
		repo := new(mock_repo.MockNotificationRepository)
		drepo := new(mock_repo.MockDeviceRepository)
		hub := new(mock_socket.MockSocketHub)
		dispatch := new(mock_dispatch.MockDispatcher)
		authClient := new(mock_client.MockAuthClient)
		lineClient := new(mock_line.MockLineClient)
		uc := usecase.NewNotificationUsecase(repo, drepo, hub, dispatch, lineClient, authClient)

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
		drepo := new(mock_repo.MockDeviceRepository)
		hub := new(mock_socket.MockSocketHub)
		dispatch := new(mock_dispatch.MockDispatcher)
		authClient := new(mock_client.MockAuthClient)
		lineClient := new(mock_line.MockLineClient)
		uc := usecase.NewNotificationUsecase(repo, drepo, hub, dispatch, lineClient, authClient)

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

		hub.On("IsOnline", targetID.String()).Return(true).Once()

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
		drepo := new(mock_repo.MockDeviceRepository)
		hub := new(mock_socket.MockSocketHub)
		dispatch := new(mock_dispatch.MockDispatcher)
		authClient := new(mock_client.MockAuthClient)
		lineClient := new(mock_line.MockLineClient)
		uc := usecase.NewNotificationUsecase(repo, drepo, hub, dispatch, lineClient, authClient)

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

		hub.On("IsOnline", targetID.String()).Return(true).Once()

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
		drepo := new(mock_repo.MockDeviceRepository)
		hub := new(mock_socket.MockSocketHub)
		dispatch := new(mock_dispatch.MockDispatcher)
		authClient := new(mock_client.MockAuthClient)
		lineClient := new(mock_line.MockLineClient)
		uc := usecase.NewNotificationUsecase(repo, drepo, hub, dispatch, lineClient, authClient)
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

		hub.On("IsOnline", targetID.String()).Return(true).Once()

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
		drepo := new(mock_repo.MockDeviceRepository)
		hub := new(mock_socket.MockSocketHub)
		dispatch := new(mock_dispatch.MockDispatcher)
		authClient := new(mock_client.MockAuthClient)
		lineClient := new(mock_line.MockLineClient)
		uc := usecase.NewNotificationUsecase(repo, drepo, hub, dispatch, lineClient, authClient)

		ctx := context.Background()

		evt := domain.InsuranceExpiringEvent{
			UserID: "not-a-uuid",
		}

		err := uc.HandleInsuranceExpiring(ctx, evt)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user id")

		repo.AssertNotCalled(t, "CreateNotification", mock.Anything, mock.Anything)
	})

	t.Run("HandleInsuranceExpiring - Success with LINE Notification", func(t *testing.T) {
		repo := new(mock_repo.MockNotificationRepository)
		drepo := new(mock_repo.MockDeviceRepository)
		hub := new(mock_socket.MockSocketHub)
		dispatch := new(mock_dispatch.MockDispatcher)
		authClient := new(mock_client.MockAuthClient)
		lineClient := new(mock_line.MockLineClient)

		uc := usecase.NewNotificationUsecase(repo, drepo, hub, dispatch, lineClient, authClient)

		ctx := context.Background()
		userID := uuid.New()
		insuranceID := uuid.New()
		lineID := "U1234567890abcdef"

		evt := domain.InsuranceExpiringEvent{
			UserID:        userID.String(),
			InsuranceID:   insuranceID.String(),
			InsuranceName: "AIA Health",
			DaysLeft:      7,
			ExpDate:       "2026-03-17",
		}

		repo.On("CreateNotification", ctx, mock.Anything).Return(nil).Once()
		hub.On("IsOnline", userID.String()).Return(false).Once()
		dispatch.On("SendToUser",
			mock.Anything,
			mock.AnythingOfType("[]domain.DeviceToken"),
			mock.AnythingOfType("push_provider.PushPayload"),
		).Return(nil).Maybe()

		drepo.On("GetActiveTokens", mock.Anything, userID).
			Return([]domain.DeviceToken{
				{
					Token: "token_1",
				},
			}, nil).
			Maybe()

		mockGrpcResponse := &pb.GetProviderAccountsResponse{
			Accounts: []*pb.ProviderAccount{
				{
					Provider:          "line",
					IsLinked:          true,
					ProviderAccountId: lineID,
				},
			},
		}
		authClient.On("GetProviderAccount", mock.Anything, mock.MatchedBy(func(req *pb.GetProviderAccountRequest) bool {
			return req.UserId == userID.String()
		})).Return(mockGrpcResponse, nil).Once()

		lineClient.On("SendTextMessage", lineID, mock.AnythingOfType("string")).Return(nil).Once()
		err := uc.HandleInsuranceExpiring(ctx, evt)
		assert.NoError(t, err)
		time.Sleep(150 * time.Millisecond)

		repo.AssertExpectations(t)
		authClient.AssertExpectations(t)
		lineClient.AssertExpectations(t)
		dispatch.AssertExpectations(t)
	})

	t.Run("HandleAccessGranted - Check Payload", func(t *testing.T) {
		repo := new(mock_repo.MockNotificationRepository)
		drepo := new(mock_repo.MockDeviceRepository)
		hub := new(mock_socket.MockSocketHub)
		dispatch := new(mock_dispatch.MockDispatcher)
		authClient := new(mock_client.MockAuthClient)
		lineClient := new(mock_line.MockLineClient)
		uc := usecase.NewNotificationUsecase(repo, drepo, hub, dispatch, lineClient, authClient)

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

		hub.On("IsOnline", targetID.String()).Return(true).Once()

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
