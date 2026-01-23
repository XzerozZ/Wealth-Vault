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
}
