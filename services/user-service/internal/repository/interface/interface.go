package repository

import (
	"context"
	"wealth-vault/user-service/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
}
