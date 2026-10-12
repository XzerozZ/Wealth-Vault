package mock

import (
	"context"
	"wealth-vault/auth-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) Register(ctx context.Context, auth *domain.AuthAccount) error {
	args := m.Called(ctx, auth)
	return args.Error(0)
}

func (m *MockAuthRepository) FindByEmailAndProvider(ctx context.Context, email string, provider string) (*domain.AuthAccount, error) {
	args := m.Called(ctx, email, provider)
	return args.Get(0).(*domain.AuthAccount), args.Error(1)
}

func (m *MockAuthRepository) FindByID(ctx context.Context, userid string) (*domain.AuthAccount, error) {
	args := m.Called(ctx, userid)
	return args.Get(0).(*domain.AuthAccount), args.Error(1)
}

func (m *MockAuthRepository) SaveOTP(ctx context.Context, otp *domain.AuthOTP) error {
	args := m.Called(ctx, otp)
	return args.Error(0)
}

func (m *MockAuthRepository) CreateSession(ctx context.Context, session *domain.AuthSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockAuthRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*domain.AuthSession, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(*domain.AuthSession), args.Error(1)
}

func (m *MockAuthRepository) RevokeSession(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockAuthRepository) RotateSession(ctx context.Context, oldRefreshToken string, next *domain.AuthSession) error {
	args := m.Called(ctx, oldRefreshToken, next)
	return args.Error(0)
}

func (m *MockAuthRepository) FindLocalAccountByEmail(ctx context.Context, email string) (*domain.AuthAccount, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AuthAccount), args.Error(1)
}

func (m *MockAuthRepository) FindLineIdentity(ctx context.Context, channelID, subject string) (*domain.LineIdentity, error) {
	args := m.Called(ctx, channelID, subject)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LineIdentity), args.Error(1)
}

func (m *MockAuthRepository) FindLineIdentityByUserID(ctx context.Context, userID uuid.UUID, channelID string) (*domain.LineIdentity, error) {
	args := m.Called(ctx, userID, channelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LineIdentity), args.Error(1)
}

func (m *MockAuthRepository) SaveLineChallenge(ctx context.Context, challenge *domain.LineChallenge) error {
	return m.Called(ctx, challenge).Error(0)
}

func (m *MockAuthRepository) ClaimLineChallenge(ctx context.Context, digest string) (*domain.LineChallenge, error) {
	args := m.Called(ctx, digest)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LineChallenge), args.Error(1)
}

func (m *MockAuthRepository) ReleaseLineChallenge(ctx context.Context, digest string) error {
	return m.Called(ctx, digest).Error(0)
}

func (m *MockAuthRepository) LinkLineIdentity(ctx context.Context, identity *domain.LineIdentity) error {
	return m.Called(ctx, identity).Error(0)
}

func (m *MockAuthRepository) DeleteExpiredLineChallenges(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *MockAuthRepository) GetValidOTP(ctx context.Context, userID uuid.UUID, code string) (*domain.AuthOTP, error) {
	args := m.Called(ctx, userID, code)
	return args.Get(0).(*domain.AuthOTP), args.Error(1)
}

func (m *MockAuthRepository) DeleteOTP(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, newHash string) error {
	args := m.Called(ctx, userID, newHash)
	return args.Error(0)
}

func (m *MockAuthRepository) DeleteExpiredSessions(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAuthRepository) DeleteExpiredOTPs(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAuthRepository) FindByUserIDAndProvider(ctx context.Context, userID uuid.UUID, provider string) (*domain.AuthAccount, error) {
	args := m.Called(ctx, userID, provider)
	return args.Get(0).(*domain.AuthAccount), args.Error(1)
}

func (m *MockAuthRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.AuthAccount, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.AuthAccount), args.Error(1)
}
