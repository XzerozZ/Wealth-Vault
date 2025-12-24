package usecase

import (
	"context"
	"wealth-vault/user-service/internal/domain"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, user *domain.User) (string, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}
