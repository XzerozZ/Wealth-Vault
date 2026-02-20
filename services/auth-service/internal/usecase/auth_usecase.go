package usecase

import (
	"context"
	"fmt"
	"time"
	"wealth-vault/auth-service/internal/domain"
	repo "wealth-vault/auth-service/internal/repository/interface"
	google "wealth-vault/auth-service/pkg/google/interface"
	mail "wealth-vault/auth-service/pkg/mail/interface"
	pb "wealth-vault/auth-service/pkg/pb/proto/user"
	token "wealth-vault/auth-service/pkg/token/interface"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	ProviderLocal  = "local"
	ProviderGoogle = "google"

	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
	TokenTypeReset   = "reset"

	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
	ResetTokenTTL   = 5 * time.Minute
	OTPTTL          = 5 * time.Minute

	SessionTTL = 1 * time.Hour
	OTPLength  = 6
)

type AuthUsecase struct {
	authRepo        repo.AuthRepository
	userClient      pb.UserServiceClient
	token           token.Generator
	mailclient      mail.NotificationClient
	googleValidator google.GoogleTokenValidator
}

func NewAuthUsecase(
	r repo.AuthRepository,
	userClient pb.UserServiceClient,
	t token.Generator,
	mail mail.NotificationClient,
	googleValidator google.GoogleTokenValidator,
) AuthUsecase {
	return AuthUsecase{
		authRepo:        r,
		userClient:      userClient,
		token:           t,
		mailclient:      mail,
		googleValidator: googleValidator,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

func (u *AuthUsecase) GenerateTokensAndSession(ctx context.Context, userID uuid.UUID, email string) (*domain.AuthOutput, error) {
	userIDStr := userID.String()

	accessToken, err := u.token.CreateToken(userIDStr, email, TokenTypeAccess, AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := u.token.CreateToken(userIDStr, email, TokenTypeRefresh, RefreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	session := &domain.AuthSession{
		UserID:           userID,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresAt:        time.Now().Add(SessionTTL),
		RefreshExpiresAt: time.Now().Add(RefreshTokenTTL),
		Revoked:          false,
	}

	if err := u.authRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &domain.AuthOutput{
		UserID:       userIDStr,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
