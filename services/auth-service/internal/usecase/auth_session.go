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

	user, err := u.authRepo.FindByID(ctx, session.UserID.String())
	if err != nil {
		return nil, errors.New("user not found")
	}

	next, output, err := u.newSession(user.UserID, user.Email)
	if err != nil {
		return nil, err
	}
	if err := u.authRepo.RotateSession(ctx, refreshTokenStr, next); err != nil {
		return nil, fmt.Errorf("failed to rotate session: %w", err)
	}
	return output, nil
}

func (u *AuthUsecase) CleanupSessions(ctx context.Context) error {
	if err := u.authRepo.DeleteExpiredSessions(ctx); err != nil {
		return fmt.Errorf("failed to cleanup sessions: %w", err)
	}
	return nil
}
