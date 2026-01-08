package repository

import (
	"context"
	"wealth-vault/auth-service/internal/domain"
)

type AuthRepository interface {
	Register(ctx context.Context, auth *domain.AuthAccount) error
	FindByEmail(ctx context.Context, email string) (*domain.AuthAccount, error)
	CreateSession(ctx context.Context, session *domain.AuthSession) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*domain.AuthSession, error)
	RevokeSession(ctx context.Context, refreshToken string) error
	DeleteExpiredSessions(ctx context.Context) error
}
