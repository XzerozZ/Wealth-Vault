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

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	authRepo   repo.AuthRepository
	userClient pb.UserServiceClient
	jwtSecret  string
}

func NewAuthUsecase(
	r repo.AuthRepository,
	userClient pb.UserServiceClient,
) AuthUsecase {
	return AuthUsecase{
		authRepo:   r,
		userClient: userClient,
	}
}

func (u *AuthUsecase) generateTokenPair(userID string, email string) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"type":    "access",
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, refreshTokenStr string) (*domain.AuthOutput, error) {
	token, err := jwt.Parse(refreshTokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return nil, errors.New("invalid token type")
	}

	session, err := u.authRepo.GetSessionByRefreshToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, errors.New("invalid or revoked refresh token")
	}

	if time.Now().After(session.RefreshExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	_ = u.authRepo.RevokeSession(ctx, refreshTokenStr)

	newAccess, newRefresh, err := u.generateTokenPair(session.UserID, "")
	if err != nil {
		return nil, err
	}

	newSession := &domain.AuthSession{
		ID:               uuid.New().String(),
		UserID:           session.UserID,
		AccessToken:      newAccess,
		RefreshToken:     newRefresh,
		ExpiresAt:        time.Now().Add(time.Hour * 1),
		RefreshExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	_ = u.authRepo.CreateSession(ctx, newSession)

	return &domain.AuthOutput{
		UserID:       session.UserID,
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
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

	accessToken, refreshToken, err := u.generateTokenPair(existingUser.ID, existingUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
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
