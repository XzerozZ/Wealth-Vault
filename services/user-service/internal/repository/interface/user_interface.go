package repository

import (
	"context"
	"wealth-vault/user-service/internal/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User, mask []string) (*domain.User, error)

	GetFriendList(ctx context.Context, userID uuid.UUID) ([]domain.FriendList, error)
	AddFriend(ctx context.Context, fri *domain.FriendList) error
	CreateFriendship(ctx context.Context, fri *domain.FriendList) error
	RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error
	UpdateFriendStatus(ctx context.Context, userID, friendID uuid.UUID, status string) error
	CheckFriendship(ctx context.Context, userID, friendID uuid.UUID) (bool, string, error)
	GetIncomingRequests(ctx context.Context, userID uuid.UUID) ([]domain.FriendList, error)
	SetCloseFriendStatus(ctx context.Context, userID, friendID uuid.UUID, isClose bool) error
	GetCloseFriends(ctx context.Context, userID uuid.UUID) ([]domain.User, error)
	GetUsersReadyForAutoShare(ctx context.Context) ([]domain.User, error)
	MarkAutoShareTriggered(ctx context.Context, userID uuid.UUID) error
	CreateFriendLog(ctx context.Context, log *domain.FriendLog) error
}
