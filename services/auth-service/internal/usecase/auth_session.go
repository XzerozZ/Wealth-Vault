package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"wealth-vault/auth-service/internal/domain"
)

func (u *AuthUsecase) RefreshToken(ctx context.Context, refreshTokenStr string) (*domain.AuthOutput, error) {
	claims, err := u.token.VerifyToken(refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != TokenTypeRefresh {
		return nil, fmt.Errorf("invalid token type")
	}

	session, err := u.authRepo.GetSessionByRefreshToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid or revoked refresh token")
	}

	if time.Now().After(session.RefreshExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	if err := u.authRepo.RevokeSession(ctx, refreshTokenStr); err != nil {
		return nil, fmt.Errorf("failed to revoke old session: %w", err)
	}

	user, err := u.authRepo.FindByID(ctx, session.UserID.String())
	if err != nil {
		return nil, errors.New("user not found")
	}

	return u.GenerateTokensAndSession(ctx, user.UserID, user.Email)
}

func (u *AuthUsecase) CleanupSessions(ctx context.Context) error {
	if err := u.authRepo.DeleteExpiredSessions(ctx); err != nil {
		return fmt.Errorf("failed to cleanup sessions: %w", err)
	}
	return nil
}
