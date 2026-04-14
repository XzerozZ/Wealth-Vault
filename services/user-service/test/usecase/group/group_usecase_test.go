package usecase_test

import (
	"context"
	"testing"
	"wealth-vault/user-service/internal/domain"
	usecase "wealth-vault/user-service/internal/usecase/group"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	mock_event "wealth-vault/user-service/test/mock/event"
	mock_repo "wealth-vault/user-service/test/mock/repository"
	mock_storage "wealth-vault/user-service/test/mock/storage"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGroupUsecase(t *testing.T) {
	ctx := context.Background()

	mockGroupRepo := new(mock_repo.MockGroupRepository)
	mockUserRepo := new(mock_repo.MockUserRepository)
	mockMsgRepo := new(mock_repo.MockMsgRepository)
	mockStorage := new(mock_storage.MockSupabaseStorage)
	mockPublisher := new(mock_event.MockEventPublisher)

	uc := usecase.NewGroupUsecase(mockGroupRepo, mockUserRepo, mockMsgRepo, mockStorage, mockPublisher)

	t.Run("CreateGroup - success", func(t *testing.T) {
		creatorID := uuid.New()
		req := &pb.CreateGroupRequest{
			Name:      "Family Group",
			CreatorId: creatorID.String(),
			MemberIds: []string{uuid.New().String()},
		}

		mockUserRepo.On("GetUser", mock.Anything, creatorID).
			Return(&domain.User{ID: creatorID, Username: "CreatorName"}, nil).Once()

		mockGroupRepo.On("CreateGroup", mock.Anything, mock.Anything, req.MemberIds).
			Return(nil).Once()

		mockPublisher.On("Publish", mock.Anything, mock.Anything).Return(nil).Once()

		resp, err := uc.CreateGroup(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Name, resp.Group.Name)
		assert.Equal(t, int64(2), resp.Group.MemberCount)
		mockGroupRepo.AssertExpectations(t)
	})

	t.Run("GetGroup - fail (not a member)", func(t *testing.T) {
		groupID := uuid.New()
		userID := uuid.New()

		mockGroupRepo.On("IsUserMember", mock.Anything, groupID, userID).
			Return(false, nil).Once()

		resp, err := uc.GetGroup(ctx, &pb.GetGroupRequest{
			GroupId: groupID.String(),
			UserId:  userID.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("GetMember - success", func(t *testing.T) {
		groupID := uuid.New()
		userID := uuid.New()

		mockGroupRepo.On("IsUserMember", mock.Anything, groupID, userID).
			Return(true, nil).Once()

		members := []*domain.User{{ID: userID, Username: "Tester"}}
		mockGroupRepo.On("GetMember", mock.Anything, groupID).
			Return(members, int64(1), nil).Once()

		resp, err := uc.GetMember(ctx, &pb.GetGroupMembersRequest{
			GroupId: groupID.String(),
			UserId:  userID.String(),
		})

		assert.NoError(t, err)
		assert.Len(t, resp.Members, 1)
		assert.Equal(t, int64(1), resp.Total)
	})

	t.Run("UpdateGroup - success", func(t *testing.T) {
		groupID := uuid.New()
		userID := uuid.New()
		existingGroup := &domain.Group{ID: groupID, GroupName: "Old Name"}

		mockGroupRepo.On("IsUserMember", mock.Anything, groupID, userID).
			Return(true, nil).Once()

		mockGroupRepo.On("GetGroupMember", mock.Anything, groupID, userID).Return(&domain.GroupMember{}, nil).Once()
		mockGroupRepo.On("GetGroup", mock.Anything, groupID).Return(existingGroup, int64(5), nil).Once()

		mockUserRepo.On("GetUser", mock.Anything, userID).Return(&domain.User{Username: "Editor"}, nil).Once()

		mockGroupRepo.On("UpdateGroup", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(&domain.Group{ID: groupID, GroupName: "New Name"}, int64(5), nil).Once()

		resp, err := uc.UpdateGroup(ctx, &pb.UpdateGroupRequest{
			Id:     groupID.String(),
			UserId: userID.String(),
			Name:   "New Name",
		})

		assert.NoError(t, err)
		assert.Equal(t, "New Name", resp.Group.Name)
	})

	t.Run("RemoveMember - success as creator", func(t *testing.T) {
		groupID := uuid.New()
		adminID := uuid.New()
		targetID := uuid.New()
		group := &domain.Group{ID: groupID, CreatedBy: adminID}

		mockGroupRepo.On("GetGroup", mock.Anything, groupID).Return(group, int64(3), nil).Once()
		mockUserRepo.On("GetUser", mock.Anything, adminID).Return(&domain.User{Username: "Admin"}, nil).Once()
		mockUserRepo.On("GetUser", mock.Anything, targetID).Return(&domain.User{Username: "Target"}, nil).Once()

		mockGroupRepo.On("RemoveMemberAndTheirSharedItems", mock.Anything, groupID, targetID, mock.Anything).
			Return(nil).Once()

		mockPublisher.On("Publish", mock.Anything, mock.Anything).Return(nil).Once()

		resp, err := uc.RemoveMember(ctx, &pb.RemoveMemberRequest{
			GroupId:        groupID.String(),
			UserId:         adminID.String(),
			TargetMemberId: targetID.String(),
		})

		assert.NoError(t, err)
		assert.True(t, resp.Success)
	})

	t.Run("LeaveGroup - fail (owner cannot leave)", func(t *testing.T) {
		groupID := uuid.New()
		ownerID := uuid.New()
		group := &domain.Group{ID: groupID, CreatedBy: ownerID}

		mockGroupRepo.On("GetGroup", mock.Anything, groupID).Return(group, int64(5), nil).Once()

		resp, err := uc.LeaveGroup(ctx, &pb.LeaveGroupRequest{
			GroupId: groupID.String(),
			UserId:  ownerID.String(),
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "owner cannot leave group")
		assert.Nil(t, resp)
	})

	t.Run("DeleteGroupCompletely - success", func(t *testing.T) {
		groupID := uuid.New()
		userID := uuid.New()
		group := &domain.Group{ID: groupID, CreatedBy: userID}

		mockGroupRepo.On("GetGroup", mock.Anything, groupID).
			Return(group, int64(1), nil).Once()

		mockGroupRepo.On("DeleteGroupCompletely", mock.Anything, groupID).
			Return(nil).Once()

		resp, err := uc.DeleteGroup(ctx, &pb.DeleteGroupRequest{
			GroupId: groupID.String(),
			UserId:  userID.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Success)

		mockGroupRepo.AssertExpectations(t)
	})
}
