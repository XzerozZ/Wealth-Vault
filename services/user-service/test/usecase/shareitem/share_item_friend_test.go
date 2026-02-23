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
	mock_repo "wealth-vault/user-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShareItemUsecase_FriendFlows(t *testing.T) {
	t.Run("GetSharedIteminFriend - success", func(t *testing.T) {
		mockRepo := new(mock_repo.MockShareItemRepository)
		mockAsset := new(mock_client.MockAssetClient)
		uc := usecase.NewShareItemUsecase(mockRepo, nil, nil, nil, mockAsset, nil, nil)

		uid := uuid.New()
		fid := uuid.New()
		itemID := uuid.New()

		items := []domain.FriendItem{
			{
				ID:         itemID,
				OwnerID:    fid,
				EntityID:   uuid.New(),
				EntityType: "ACCOUNT",
				CreatedAt:  time.Now(),
			},
		}

		mockRepo.On("GetSharedIteminFriend", mock.Anything, fid, uid).
			Return(items, nil)

		mockAsset.On("GetBatchAccount", mock.Anything, mock.Anything, mock.Anything).
			Return(&assetPb.AccountArrayResponse{
				Account: []*assetPb.Account{
					{Id: itemID.String(), Name: "Test Account"},
				},
			}, nil)

		resp, err := uc.GetSharedIteminFriend(context.Background(), &pb.GetFriendItemRequest{
			UserId:   uid.String(),
			FriendId: fid.String(),
		})

		assert.NoError(t, err)
		assert.Len(t, resp.Items, 1)
	})

	t.Run("UnsharedIteminFriend - success", func(t *testing.T) {
		mockRepo := new(mock_repo.MockShareItemRepository)

		uc := usecase.NewShareItemUsecase(mockRepo, nil, nil, nil, nil, nil, nil)

		itemID := uuid.New()
		userID := uuid.New()

		mockRepo.On("DeleteIteminFriend", mock.Anything, itemID, userID).
			Return(nil)

		resp, err := uc.UnsharedIteminFriend(context.Background(), &pb.UnshareItemRequest{
			ItemId: itemID.String(),
			UserId: userID.String(),
		})

		assert.NoError(t, err)
		assert.True(t, resp.Finish)
	})

	t.Run("BatchShareAssets - insert new items", func(t *testing.T) {
		mockRepo := new(mock_repo.MockShareItemRepository)
		mockUser := new(mock_repo.MockUserRepository)
		mockPub := new(mock_event.MockEventPublisher)

		uc := usecase.NewShareItemUsecase(mockRepo, nil, mockUser, nil, nil, nil, mockPub)

		ownerID := uuid.New()
		targetID := uuid.New()
		accountID := uuid.New()

		existing := map[string]bool{}

		mockRepo.On("GetExistingSharedMap", mock.Anything, ownerID, targetID).
			Return(existing, nil)

		mockRepo.On("ShareItemtoFriend", mock.Anything, mock.Anything).
			Return(nil)

		mockUser.On("GetUser", mock.Anything, ownerID).
			Return(&domain.User{Username: "Tester"}, nil)

		mockPub.On("Publish", "noti.item.shared", mock.Anything).
			Return(nil)

		err := uc.BatchShareAssets(context.Background(), domain.BatchShareRequest{
			OwnerID:    ownerID,
			TargetID:   targetID,
			AccountIDs: []string{accountID.String()},
		})

		assert.NoError(t, err)
	})

	t.Run("GetItemsSharedByFriend - success", func(t *testing.T) {
		mockRepo := new(mock_repo.MockShareItemRepository)
		mockAsset := new(mock_client.MockAssetClient)

		uc := usecase.NewShareItemUsecase(mockRepo, nil, nil, nil, mockAsset, nil, nil)
		uid := uuid.New()
		fid := uuid.New()

		summaries := []domain.SharedItemSummary{
			{EntityID: "123", EntityType: "ACCOUNT"},
		}

		mockRepo.On("GetItemsSharedByFriend", mock.Anything, uid, fid).
			Return(summaries, nil)

		mockAsset.On("GetBatchAccount", mock.Anything, mock.Anything, mock.Anything).
			Return(&assetPb.AccountArrayResponse{
				Account: []*assetPb.Account{{Id: "123"}},
			}, nil)

		resp, err := uc.GetItemsSharedByFriend(context.Background(), &pb.GetItemsSharedByFriendRequest{
			UserId:   uid.String(),
			FriendId: fid.String(),
		})

		assert.NoError(t, err)
		assert.Len(t, resp.AssetDetail, 1)
	})
}
