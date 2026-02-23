package usecase_test

import (
	"context"
	"errors"
	"testing"

	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/usecase"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	mock_repo "wealth-vault/user-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetGroupMessages(t *testing.T) {
	mockItemRepo := new(mock_repo.MockShareItemRepository)
	mockMsgRepo := new(mock_repo.MockMsgRepository)
	uc := usecase.NewMessageUsecase(mockMsgRepo, mockItemRepo)

	ctx := context.Background()
	validGID := uuid.New()
	validUID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockItemRepo.On("IsGroupMember", ctx, validGID, validUID).Return(true, nil).Once()

		mockMsgs := []domain.GroupMessage{{ID: uuid.New()}}
		mockMsgRepo.On("GetGroupMessages", ctx, validGID.String()).Return(mockMsgs, nil).Once()

		req := &pb.GetGroupMessagesRequest{
			GroupId: validGID.String(),
			UserId:  validUID.String(),
		}
		res, err := uc.GetGroupMessages(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.Messages, 1)
		mockItemRepo.AssertExpectations(t)
		mockMsgRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized - Not a member", func(t *testing.T) {
		mockItemRepo.On("IsGroupMember", ctx, validGID, validUID).Return(false, nil).Once()

		req := &pb.GetGroupMessagesRequest{
			GroupId: validGID.String(),
			UserId:  validUID.String(),
		}
		res, err := uc.GetGroupMessages(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "unauthorized", err.Error())
		assert.Nil(t, res)
	})
}

func TestGetPrivateMessages(t *testing.T) {
	mockItemRepo := new(mock_repo.MockShareItemRepository)
	mockMsgRepo := new(mock_repo.MockMsgRepository)
	uc := usecase.NewMessageUsecase(mockMsgRepo, mockItemRepo)

	ctx := context.Background()
	uid := "user-1"
	fid := "friend-1"

	t.Run("Repository Error", func(t *testing.T) {
		var nilMsgs []domain.PrivateMessage
		mockMsgRepo.On("GetPrivateMessages", ctx, uid, fid).Return(nilMsgs, errors.New("db error")).Once()

		req := &pb.GetPrivateMessagesRequest{UserId: uid, FriendId: fid}
		res, err := uc.GetPrivateMessages(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		assert.Nil(t, res)
	})
}
