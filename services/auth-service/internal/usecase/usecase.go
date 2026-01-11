package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"wealth-vault/auth-service/internal/domain"
	repo "wealth-vault/auth-service/internal/repository/interface"
	pb "wealth-vault/auth-service/pkg/pb/proto/user"
	"wealth-vault/auth-service/pkg/token"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	authRepo   repo.AuthRepository
	userClient pb.UserServiceClient
	token      token.Generator
}

func NewAuthUsecase(
	r repo.AuthRepository,
	userClient pb.UserServiceClient,
	t token.Generator,
) AuthUsecase {
	return AuthUsecase{
		authRepo:   r,
		userClient: userClient,
		token:      t,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, input *domain.RegisterInput) (*domain.AuthOutput, error) {
	existingUser, err := u.authRepo.FindByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	username := input.Username
	if username == "" {
		username = strings.Split(input.Email, "@")[0]
	}

	userRes, err := u.userClient.CreateUser(ctx, &pb.CreateUserRequest{
		Email:    input.Email,
		Username: username,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user profile: %v", err)
	}

	hashedPwd := ""
	if input.Provider == "local" {
		bytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashedPwd = string(bytes)
	}

	newAuth := &domain.AuthAccount{
		UserID:          userRes.Id,
		Email:           input.Email,
		Password:        hashedPwd,
		Provider:        input.Provider,
		IsEmailVerified: false,
	}

	if err := u.authRepo.Register(ctx, newAuth); err != nil {
		return nil, fmt.Errorf("failed to register auth account: %w", err)
	}

	return &domain.AuthOutput{
		UserID:       userRes.Id,
		AccessToken:  "",
		RefreshToken: "",
	}, nil
}

func (u *AuthUsecase) Login(ctx context.Context, input *domain.LoginInput) (*domain.AuthOutput, error) {
	existingUser, err := u.authRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(input.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	accessToken, err := u.token.CreateToken(existingUser.UserID, existingUser.Email, "access", 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := u.token.CreateToken(existingUser.UserID, existingUser.Email, "refresh", 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	session := &domain.AuthSession{
		ID:               uuid.New().String(),
		UserID:           existingUser.UserID,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresAt:        time.Now().Add(time.Hour * 1),
		RefreshExpiresAt: time.Now().Add(time.Hour * 24 * 7),
		Revoked:          false,
	}

	if err := u.authRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &domain.AuthOutput{
		UserID:       existingUser.UserID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, refreshTokenStr string) (*domain.AuthOutput, error) {
	claims, err := u.token.VerifyToken(refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
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
		return nil, err
	}

	user, err := u.authRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	newAccess, err := u.token.CreateToken(user.UserID, user.Email, "access", 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	newRefresh, err := u.token.CreateToken(user.UserID, user.Email, "refresh", 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	newSession := &domain.AuthSession{
		ID:               uuid.New().String(),
		UserID:           session.UserID,
		AccessToken:      newAccess,
		RefreshToken:     newRefresh,
		ExpiresAt:        time.Now().Add(time.Hour * 1),
		RefreshExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	if err = u.authRepo.CreateSession(ctx, newSession); err != nil {
		return nil, err
	}

	return &domain.AuthOutput{
		UserID:       session.UserID,
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}

func (u *AuthUsecase) CleanupSessions(ctx context.Context) error {
	if err := u.authRepo.DeleteExpiredSessions(ctx); err != nil {
		return fmt.Errorf("failed to cleanup sessions: %w", err)
	}

	return nil
}
