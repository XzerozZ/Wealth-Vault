package usecase_test

import (
	"context"
	"testing"
	"wealth-vault/user-service/internal/domain"
	usecase "wealth-vault/user-service/internal/usecase/user"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	mock_event "wealth-vault/user-service/test/mock/event"
	mock_repo "wealth-vault/user-service/test/mock/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddFriend(t *testing.T) {
	ctx := context.Background()

	t.Run("cannot add yourself", func(t *testing.T) {
		userRepo := new(mock_repo.MockUserRepository)
		pub := new(mock_event.MockEventPublisher)
		uc := usecase.NewUserUsecase(userRepo, nil, nil, pub, nil)

		id := uuid.New().String()

		_, err := uc.AddFriend(ctx, &pb.FriendRequest{
			Id:     id,
			UserId: id,
		})

		assert.Error(t, err)
		assert.Equal(t, "cannot add yourself", err.Error())
	})

	t.Run("already friends", func(t *testing.T) {
		userRepo := new(mock_repo.MockUserRepository)
		pub := new(mock_event.MockEventPublisher)
		uc := usecase.NewUserUsecase(userRepo, nil, nil, pub, nil)

		u1 := uuid.New()
		u2 := uuid.New()

		userRepo.
			On("CheckFriendship", ctx, u1, u2).
			Return(true, "ACCEPTED", nil)

		_, err := uc.AddFriend(ctx, &pb.FriendRequest{
			Id:     u1.String(),
			UserId: u2.String(),
		})

		assert.Error(t, err)
		assert.Equal(t, "already friends", err.Error())
	})

	t.Run("reverse pending", func(t *testing.T) {
		userRepo := new(mock_repo.MockUserRepository)
		pub := new(mock_event.MockEventPublisher)
		uc := usecase.NewUserUsecase(userRepo, nil, nil, pub, nil)

		u1 := uuid.New()
		u2 := uuid.New()

		userRepo.
			On("CheckFriendship", ctx, u1, u2).
			Return(false, "", nil)

		userRepo.
			On("CheckFriendship", ctx, u2, u1).
			Return(true, "PENDING", nil)

		_, err := uc.AddFriend(ctx, &pb.FriendRequest{
			Id:     u1.String(),
			UserId: u2.String(),
		})

		assert.Error(t, err)
		assert.Equal(t, "this user has already sent you a friend request, please check your pending requests", err.Error())
	})

	t.Run("success", func(t *testing.T) {
		userRepo := new(mock_repo.MockUserRepository)
		pub := new(mock_event.MockEventPublisher)
		uc := usecase.NewUserUsecase(userRepo, nil, nil, pub, nil)

		u1 := uuid.New()
		u2 := uuid.New()

		userRepo.
			On("CheckFriendship", ctx, u1, u2).
			Return(false, "", nil)

		userRepo.
			On("CheckFriendship", ctx, u2, u1).
			Return(false, "", nil)

		userRepo.
			On("AddFriend", ctx, mock.Anything).
			Return(nil)

		userRepo.
			On("GetUser", ctx, u1).
			Return(&domain.User{Username: "John"}, nil)

		pub.
			On("Publish", "noti.friend.request", mock.Anything).
			Return(nil)

		res, err := uc.AddFriend(ctx, &pb.FriendRequest{
			Id:     u1.String(),
			UserId: u2.String(),
		})

		assert.NoError(t, err)
		assert.True(t, res.Success)
	})
}

func TestAcceptFriend(t *testing.T) {
	ctx := context.Background()

	t.Run("decline", func(t *testing.T) {
		userRepo := new(mock_repo.MockUserRepository)
		pub := new(mock_event.MockEventPublisher)
		uc := usecase.NewUserUsecase(userRepo, nil, nil, pub, nil)

		u1 := uuid.New()
		u2 := uuid.New()

		userRepo.
			On("RemoveFriend", ctx, u1, u2).
			Return(nil)

		res, err := uc.AcceptFriend(ctx, &pb.AcceptFriendRequest{
			UserId:      u1.String(),
			RequesterId: u2.String(),
			Action:      "DECLINE",
		})

		assert.NoError(t, err)
		assert.True(t, res.Success)
	})

	t.Run("accept", func(t *testing.T) {
		userRepo := new(mock_repo.MockUserRepository)
		pub := new(mock_event.MockEventPublisher)
		uc := usecase.NewUserUsecase(userRepo, nil, nil, pub, nil)

		u1 := uuid.New()
		u2 := uuid.New()

		userRepo.
			On("UpdateFriendStatus", ctx, u1, u2, "ACCEPTED").
			Return(nil)

		userRepo.
			On("CreateFriendship", ctx, mock.Anything).
			Return(nil)

		userRepo.
			On("GetUser", ctx, u1).
			Return(&domain.User{Username: "Alice"}, nil)

		pub.
			On("Publish", "noti.friend.accepted", mock.Anything).
			Return(nil)

		res, err := uc.AcceptFriend(ctx, &pb.AcceptFriendRequest{
			UserId:      u1.String(),
			RequesterId: u2.String(),
		})

		assert.NoError(t, err)
		assert.True(t, res.Success)
	})
}

func TestSetCloseFriend(t *testing.T) {
	ctx := context.Background()

	userRepo := new(mock_repo.MockUserRepository)
	pub := new(mock_event.MockEventPublisher)
	uc := usecase.NewUserUsecase(userRepo, nil, nil, pub, nil)

	u1 := uuid.New()
	u2 := uuid.New()

	userRepo.
		On("CheckFriendship", ctx, u1, u2).
		Return(true, "ACCEPTED", nil)

	userRepo.
		On("SetCloseFriendStatus", ctx, u1, u2, true).
		Return(nil)

	res, err := uc.SetCloseFriend(ctx, &pb.SetCloseFriendRequest{
		UserId:   u1.String(),
		FriendId: u2.String(),
		IsClose:  true,
	})

	assert.NoError(t, err)
	assert.True(t, res.Success)
	userRepo.AssertExpectations(t)
}

func TestGetCloseFriends_Usecase(t *testing.T) {
	ctx := context.Background()
	userRepo := new(mock_repo.MockUserRepository)
	pub := new(mock_event.MockEventPublisher)

	uc := usecase.NewUserUsecase(userRepo, nil, nil, pub, nil)

	u1 := uuid.New()
	friendID := uuid.New()
	userRepo.
		On("GetCloseFriends", mock.Anything, u1).
		Return([]domain.FriendList{
			{
				UserID:        u1,
				FriendID:      friendID,
				Status:        "ACCEPTED",
				IsCloseFriend: true,
				Friend: domain.User{
					ID:       friendID,
					Username: "Bob",
				},
			},
		}, nil)

	res, err := uc.GetCloseFriends(ctx, &pb.GetCloseFriendsRequest{
		UserId: u1.String(),
	})

	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Friends, 1)

	assert.Equal(t, "Bob", res.Friends[0].Username)
	assert.True(t, res.Friends[0].IsCloseFriend)

	userRepo.AssertExpectations(t)
}
