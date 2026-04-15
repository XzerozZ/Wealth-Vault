package usecase_test

import (
	"context"
	"testing"
	"time"
	"wealth-vault/user-service/internal/domain"
	usecase "wealth-vault/user-service/internal/usecase/shareItem"
	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	mock_client "wealth-vault/user-service/test/mock/client"
	mock_event "wealth-vault/user-service/test/mock/event"
	mock_mail "wealth-vault/user-service/test/mock/mail"
	mock_repo "wealth-vault/user-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShareItemUsecase_ShareItem(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	itemID := uuid.New()
	groupID := uuid.New()
	friendID := uuid.New()

	t.Run("❌ no items", func(t *testing.T) {
		uc := usecase.NewShareItemUsecase(nil, nil, nil, nil, nil, nil, nil)
		resp, err := uc.ShareItem(ctx, &pb.ShareItemRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("❌ mismatch itemIds and types", func(t *testing.T) {
		uc := usecase.NewShareItemUsecase(nil, nil, nil, nil, nil, nil, nil)
		resp, err := uc.ShareItem(ctx, &pb.ShareItemRequest{
			ItemIds:   []string{"1"},
			ItemTypes: []string{},
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("❌ asset not found", func(t *testing.T) {
		mockItem := new(mock_repo.MockShareItemRepository)
		mockUser := new(mock_repo.MockUserRepository)
		mockAsset := new(mock_client.MockAssetClient)
		mockPub := new(mock_event.MockEventPublisher)

		uc := usecase.NewShareItemUsecase(mockItem, nil, mockUser, nil, mockAsset, nil, mockPub)

		mockUser.On("GetUser", mock.Anything, mock.Anything).
			Return(&domain.User{Username: "Tester"}, nil)

		mockAsset.On("CheckAssetExists", mock.Anything, mock.Anything, mock.Anything).
			Return(&assetPb.CheckAssetResponse{Exists: false}, nil)

		resp, err := uc.ShareItem(ctx, &pb.ShareItemRequest{
			UserId:    userID.String(),
			ItemIds:   []string{itemID.String()},
			ItemTypes: []string{"ASSET"},
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("✅ success share group + friend + email", func(t *testing.T) {
		mockItem := new(mock_repo.MockShareItemRepository)
		mockGroup := new(mock_repo.MockGroupRepository)
		mockUser := new(mock_repo.MockUserRepository)
		mockMsg := new(mock_repo.MockMsgRepository)
		mockAsset := new(mock_client.MockAssetClient)
		mockMail := new(mock_mail.MockMailClient)
		mockPub := new(mock_event.MockEventPublisher)
		uc := usecase.NewShareItemUsecase(mockItem, mockGroup, mockUser, mockMsg, mockAsset, mockMail, mockPub)

		mockUser.On("GetUser", mock.Anything, userID).
			Return(&domain.User{Username: "Tester"}, nil)

		mockMsg.On(
			"CreatePrivateMessage",
			mock.Anything,
			mock.AnythingOfType("[]domain.PrivateMessage"),
		).Return(nil)

		mockUser.On(
			"CreateFriendLog",
			mock.Anything,
			mock.AnythingOfType("*domain.FriendLog"),
		).Return(nil)

		mockMail.On(
			"SendShareInvitation",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil)

		mockAsset.On("CheckAssetExists", mock.Anything, mock.Anything, mock.Anything).
			Return(&assetPb.CheckAssetResponse{Exists: true}, nil)

		mockItem.On("IsItemSharedtoGroup", mock.Anything, groupID, itemID, "ASSET").
			Return(false, nil)

		mockItem.On("IsItemSharedtoFriend", mock.Anything, friendID, itemID, "ASSET").
			Return(false, nil)

		mockGroup.On("GetMember", mock.Anything, groupID).
			Return([]*domain.User{
				{
					ID:       friendID,
					Username: "Friend",
				},
			}, int64(1), nil)

		mockItem.On("ShareItemtoGroup", mock.Anything, mock.Anything).
			Return(nil)

		mockItem.On("ShareItemtoFriend", mock.Anything, mock.Anything).
			Return(nil)

		mockItem.On("ShareItemtoEmail", mock.Anything, mock.Anything).
			Return(nil)

		mockGroup.On(
			"CreateLog",
			mock.Anything,
			mock.AnythingOfType("*domain.GroupLog"),
		).Return(nil).Once()

		mockMsg.On(
			"CreateMessage",
			mock.Anything,
			mock.AnythingOfType("[]domain.GroupMessage"),
		).Return(nil)

		mockPub.On("Publish", mock.Anything, mock.Anything).
			Return(nil)

		resp, err := uc.ShareItem(ctx, &pb.ShareItemRequest{
			UserId:    userID.String(),
			ItemIds:   []string{itemID.String()},
			ItemTypes: []string{"ASSET"},
			Groups: []*pb.ShareTarget{
				{Id: groupID.String()},
			},
			Friends: []*pb.ShareTarget{
				{Id: friendID.String()},
			},
			Emails: []*pb.ShareTarget{
				{Id: "test@email.com"},
			},
		})

		assert.NoError(t, err)
		assert.True(t, resp.Finish)
		assert.Eventually(t, func() bool {
			return mockMail.AssertExpectations(t)
		}, 1*time.Second, 50*time.Millisecond)

		mockItem.AssertExpectations(t)
		mockGroup.AssertExpectations(t)
		mockUser.AssertExpectations(t)
		mockMsg.AssertExpectations(t)
		mockAsset.AssertExpectations(t)
		mockPub.AssertExpectations(t)
	})
}
