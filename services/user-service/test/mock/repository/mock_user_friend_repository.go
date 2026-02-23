package mock

import (
	"context"
	"wealth-vault/user-service/internal/domain"

	"github.com/google/uuid"
)

func (m *MockUserRepository) GetFriendList(ctx context.Context, userID uuid.UUID) ([]domain.FriendList, error) {
	args := m.Called(ctx, userID)

	var list []domain.FriendList
	if args.Get(0) != nil {
		list = args.Get(0).([]domain.FriendList)
	}

	return list, args.Error(1)
}

func (m *MockUserRepository) AddFriend(ctx context.Context, fri *domain.FriendList) error {
	args := m.Called(ctx, fri)
	return args.Error(0)
}

func (m *MockUserRepository) CreateFriendship(ctx context.Context, fri *domain.FriendList) error {
	args := m.Called(ctx, fri)
	return args.Error(0)
}

func (m *MockUserRepository) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	args := m.Called(ctx, userID, friendID)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateFriendStatus(ctx context.Context, userID, friendID uuid.UUID, status string) error {
	args := m.Called(ctx, userID, friendID, status)
	return args.Error(0)
}

func (m *MockUserRepository) CheckFriendship(ctx context.Context, userID, friendID uuid.UUID) (bool, string, error) {
	args := m.Called(ctx, userID, friendID)

	return args.Bool(0), args.String(1), args.Error(2)
}

func (m *MockUserRepository) GetIncomingRequests(ctx context.Context, userID uuid.UUID) ([]domain.FriendList, error) {
	args := m.Called(ctx, userID)

	var list []domain.FriendList
	if args.Get(0) != nil {
		list = args.Get(0).([]domain.FriendList)
	}

	return list, args.Error(1)
}

func (m *MockUserRepository) SetCloseFriendStatus(ctx context.Context, userID, friendID uuid.UUID, isClose bool) error {
	args := m.Called(ctx, userID, friendID, isClose)
	return args.Error(0)
}

func (m *MockUserRepository) GetCloseFriends(ctx context.Context, userID uuid.UUID) ([]domain.User, error) {
	args := m.Called(ctx, userID)

	var users []domain.User
	if args.Get(0) != nil {
		users = args.Get(0).([]domain.User)
	}

	return users, args.Error(1)
}

func (m *MockUserRepository) GetUsersReadyForAutoShare(ctx context.Context) ([]domain.User, error) {
	args := m.Called(ctx)

	var users []domain.User
	if args.Get(0) != nil {
		users = args.Get(0).([]domain.User)
	}

	return users, args.Error(1)
}

func (m *MockUserRepository) MarkAutoShareTriggered(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) CreateFriendLog(ctx context.Context, log *domain.FriendLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}
