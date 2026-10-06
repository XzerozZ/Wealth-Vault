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
	RotateSession(ctx context.Context, oldRefreshToken string, next *domain.AuthSession) error
	GetValidOTP(ctx context.Context, userID uuid.UUID, code string) (*domain.AuthOTP, error)
	DeleteOTP(ctx context.Context, userID uuid.UUID) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, newHash string) error
	DeleteExpiredSessions(ctx context.Context) error
	DeleteExpiredOTPs(ctx context.Context) error

	FindByUserIDAndProvider(ctx context.Context, userID uuid.UUID, provider string) (*domain.AuthAccount, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.AuthAccount, error)
	FindLocalAccountByEmail(ctx context.Context, email string) (*domain.AuthAccount, error)
	FindLineIdentity(ctx context.Context, channelID, subject string) (*domain.LineIdentity, error)
	SaveLineChallenge(ctx context.Context, challenge *domain.LineChallenge) error
	ClaimLineChallenge(ctx context.Context, digest string) (*domain.LineChallenge, error)
	ReleaseLineChallenge(ctx context.Context, digest string) error
	LinkLineIdentity(ctx context.Context, identity *domain.LineIdentity) error
	DeleteExpiredLineChallenges(ctx context.Context) error
}
