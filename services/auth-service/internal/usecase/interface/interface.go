package usecase

import (
	"context"
	"wealth-vault/auth-service/internal/domain"
)

type AuthUsecase interface {
	Register(ctx context.Context, input *domain.RegisterInput) (*domain.AuthOutput, error)
	Login(ctx context.Context, input *domain.LoginInput) (*domain.AuthOutput, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthOutput, error)
	CleanupSessions(ctx context.Context) error
}
