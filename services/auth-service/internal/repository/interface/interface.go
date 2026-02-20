package repository

import (
	"context"
	"wealth-vault/auth-service/internal/domain"

	"github.com/google/uuid"
)

type AuthRepository interface {
	Register(ctx context.Context, auth *domain.AuthAccount) error
	FindByEmailAndProvider(ctx context.Context, email string, provider string) (*domain.AuthAccount, error)
	FindByID(ctx context.Context, userid string) (*domain.AuthAccount, error)
	SaveOTP(ctx context.Context, otp *domain.AuthOTP) error
	CreateSession(ctx context.Context, session *domain.AuthSession) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*domain.AuthSession, error)
	RevokeSession(ctx context.Context, refreshToken string) error
	GetValidOTP(ctx context.Context, userID uuid.UUID, code string) (*domain.AuthOTP, error)
	DeleteOTP(ctx context.Context, userID uuid.UUID) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, newHash string) error
	DeleteExpiredSessions(ctx context.Context) error
	DeleteExpiredOTPs(ctx context.Context) error
}
