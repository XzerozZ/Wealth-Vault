package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"
	"wealth-vault/user-service/internal/domain"
	usecase "wealth-vault/user-service/internal/usecase/shareItem"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	mock_repo "wealth-vault/user-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAllReferencesByEntityID(t *testing.T) {
	ctx := context.Background()

	t.Run("invalid uuid", func(t *testing.T) {
		uc := usecase.NewShareItemUsecase(nil, nil, nil, nil, nil, nil, nil)
		_, err := uc.DeleteAllReferencesByEntityID(ctx, &pb.DeleteByEntityRequest{
			EntityId: "invalid",
		})
		assert.Error(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(mock_repo.MockShareItemRepository)
		uc := usecase.NewShareItemUsecase(repo, nil, nil, nil, nil, nil, nil)

		id := uuid.New()

		repo.
			On("DeleteAllReferencesByEntityID", ctx, id).
			Return(errors.New("db error"))

		_, err := uc.DeleteAllReferencesByEntityID(ctx, &pb.DeleteByEntityRequest{
			EntityId: id.String(),
		})

		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		repo := new(mock_repo.MockShareItemRepository)
		uc := usecase.NewShareItemUsecase(repo, nil, nil, nil, nil, nil, nil)

		id := uuid.New()

		repo.
			On("DeleteAllReferencesByEntityID", ctx, id).
			Return(nil)

		res, err := uc.DeleteAllReferencesByEntityID(ctx, &pb.DeleteByEntityRequest{
			EntityId: id.String(),
		})

		assert.NoError(t, err)
		assert.True(t, res.Success)
	})
}

func TestGetItemSharedTargets(t *testing.T) {
	ctx := context.Background()

	t.Run("success mapping", func(t *testing.T) {
		repo := new(mock_repo.MockShareItemRepository)
		uc := usecase.NewShareItemUsecase(repo, nil, nil, nil, nil, nil, nil)

		userID := uuid.New()
		itemID := uuid.New()
		now := time.Now()

		mockResp := &domain.SharedTargetsResult{
			Groups: []domain.SharedGroupTarget{
				{
					GroupID:     "g1",
					GroupName:   "Group1",
					GroupImage:  "img",
					MemberCount: 5,
					SharedAt:    now,
				},
			},
			Friends: []domain.SharedFriendTarget{
				{
					FriendID:     "f1",
					Username:     "John",
					ProfileImage: "pic",
					SharedAt:     now,
				},
			},
			Emails: []domain.SharedEmailTarget{
				{
					Email:    "a@test.com",
					SharedAt: now,
					IsSent:   true,
				},
			},
		}

		repo.
			On("GetItemSharedTargets", ctx, userID, itemID, "ACCOUNT").
			Return(mockResp, nil)

		res, err := uc.GetItemSharedTargets(ctx, &pb.GetItemSharedTargetsRequest{
			UserId:   userID.String(),
			ItemId:   itemID.String(),
			ItemType: "ACCOUNT",
		})

		assert.NoError(t, err)
		assert.Len(t, res.Groups, 1)
		assert.Len(t, res.Friends, 1)
		assert.Len(t, res.Emails, 1)

		assert.Equal(t, "g1", res.Groups[0].GroupId)
		assert.True(t, res.Emails[0].IsSent)
	})
}

func TestGetSharedItemIDs(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mock_repo.MockShareItemRepository)
		uc := usecase.NewShareItemUsecase(repo, nil, nil, nil, nil, nil, nil)

		userID := uuid.New()
		targetID := uuid.New()

		repo.
			On("GetSharedItemIDs", ctx, userID, targetID, "FRIEND").
			Return([]string{"id1", "id2"}, nil)

		res, err := uc.GetSharedItemIDs(ctx, &pb.GetSharedItemIDsRequest{
			UserId:     userID.String(),
			TargetId:   targetID.String(),
			TargetType: "FRIEND",
		})

		assert.NoError(t, err)
		assert.Equal(t, 2, len(res.ItemIds))
	})
}
