package repository

import (
	"context"
	"wealth-vault/user-service/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User, mask []string) (*domain.User, error)
}
