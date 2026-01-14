package usecase

import (
	"context"
	"wealth-vault/user-service/internal/domain"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, user *domain.User) (string, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, input *domain.UpdateUserInput) (*domain.User, error)
}
