package usecase

import (
	"context"
	"wealth-vault/auth-service/internal/domain"
)

type AuthUsecase interface {
	Register(ctx context.Context, auth *domain.AuthAccount) (*domain.AuthAccount, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthOutput, error)
}
