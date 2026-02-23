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

func TestShareItemUsecase_GroupFunctions(t *testing.T) {
	t.Run("GetSharedIteminGroup success", func(t *testing.T) {
		mockItem := new(mock_repo.MockShareItemRepository)
		mockAsset := new(mock_client.MockAssetClient)
		uc := usecase.NewShareItemUsecase(mockItem, nil, nil, nil, mockAsset, nil, nil)

		gid := uuid.New()
		uid := uuid.New()
		entityID := uuid.New()

		items := []domain.GroupItem{
			{
				ID:         uuid.New(),
				EntityID:   entityID,
				EntityType: "building",
				OwnerID:    uid,
				CreatedAt:  time.Now(),
			},
		}

		mockAsset.On("GetBatchBuilding", mock.Anything, mock.Anything, mock.Anything).
			Return(&assetPb.BuildingArrayResponse{
				Building: []*assetPb.Building{
					{
						Id:       entityID.String(),
						Name:     "ตึกทดสอบ",
						Location: &assetPb.Location{District: "เขต", Province: "จังหวัด"},
					},
				},
			}, nil)

		mockItem.On("GetSharedIteminGroup", mock.Anything, gid, uid).
			Return(items, nil)

		res, err := uc.GetSharedIteminGroup(context.Background(), &pb.GetGroupItemsRequest{
			GroupId: gid.String(),
			UserId:  uid.String(),
		})

		assert.NoError(t, err)
		assert.Len(t, res.Items, 1)
	})

	t.Run("UnsharedIteminGroup success", func(t *testing.T) {
		mockItem := new(mock_repo.MockShareItemRepository)

		uc := usecase.NewShareItemUsecase(mockItem, nil, nil, nil, nil, nil, nil)

		itemID := uuid.New()
		userID := uuid.New()

		mockItem.On("DeleteIteminGroup", mock.Anything, itemID, userID).
			Return(nil)

		res, err := uc.UnsharedIteminGroup(context.Background(), &pb.UnshareItemRequest{
			ItemId: itemID.String(),
			UserId: userID.String(),
		})

		assert.NoError(t, err)
		assert.True(t, res.Finish)
	})

	t.Run("AddMemberToGroup success", func(t *testing.T) {
		mockItem := new(mock_repo.MockShareItemRepository)
		mockGroup := new(mock_repo.MockGroupRepository)
		mockUser := new(mock_repo.MockUserRepository)
		mockMsg := new(mock_repo.MockMsgRepository)
		mockPub := new(mock_event.MockEventPublisher)

		uc := usecase.NewShareItemUsecase(mockItem, mockGroup, mockUser, mockMsg, nil, nil, mockPub)

		gid := uuid.New()
		senderID := uuid.New()
		targetID := uuid.New()

		mockUser.On("GetUser", mock.Anything, senderID).
			Return(&domain.User{Username: "Owner"}, nil)

		mockUser.On("GetUser", mock.Anything, targetID).
			Return(&domain.User{Username: "Target"}, nil)

		mockItem.On("AddMember", mock.Anything, mock.Anything).
			Return(nil)

		mockGroup.On("CreateLog", mock.Anything, mock.Anything).
			Return(nil)

		mockMsg.On("CreateMessage", mock.Anything, mock.Anything).
			Return(nil)

		mockPub.On("Publish", mock.Anything, mock.Anything).
			Return(nil)

		res, err := uc.AddMemberToGroup(context.Background(), &pb.AddMemberRequest{
			GroupId:       gid.String(),
			UserId:        senderID.String(),
			TargetUserIds: []string{targetID.String()},
		})

		assert.NoError(t, err)
		assert.True(t, res.Success)

		time.Sleep(100 * time.Millisecond)
	})

	t.Run("GrantAccess success", func(t *testing.T) {
		mockItem := new(mock_repo.MockShareItemRepository)

		uc := usecase.NewShareItemUsecase(mockItem, nil, nil, nil, nil, nil, nil)

		groupID := uuid.New()
		ownerID := uuid.New()
		targetID := uuid.New()
		itemID := uuid.New()

		mockItem.On("IsGroupMember", mock.Anything, groupID, targetID).
			Return(true, nil)

		mockItem.On("GetOwnedItemIDs", mock.Anything, []string{itemID.String()}, ownerID).
			Return([]uuid.UUID{itemID}, nil)

		mockItem.On("BatchCreateViewers", mock.Anything, mock.Anything).
			Return(nil)

		res, err := uc.GrantAccess(context.Background(), &pb.GrantAccessRequest{
			GroupId:      groupID.String(),
			OwnerUserId:  ownerID.String(),
			TargetUserId: targetID.String(),
			GroupItemIds: []string{itemID.String()},
		})

		assert.NoError(t, err)
		assert.True(t, res.Success)
	})
}
