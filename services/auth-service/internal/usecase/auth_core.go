package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"wealth-vault/auth-service/internal/domain"
	pb "wealth-vault/auth-service/pkg/pb/proto/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (u *AuthUsecase) Register(ctx context.Context, input *domain.RegisterInput) (*domain.AuthOutput, error) {
	existingUser, err := u.authRepo.FindByEmailAndProvider(ctx, input.Email, ProviderLocal)
	if existingUser != nil {
		return nil, errors.New("email already registered with this provider")
	}

	username := input.Username
	if username == "" {
		username = strings.Split(input.Email, "@")[0]
	}

	if input.Password == "" {
		return nil, errors.New("password is required for local registration")
	}

	hashedPwd, err := HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	userRes, err := u.userClient.CreateUser(ctx, &pb.CreateUserRequest{
		Email:    input.Email,
		Username: username,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user profile via gRPC: %v", err)
	}

	userID, err := uuid.Parse(userRes.User.Id)
	if err != nil {
		return nil, errors.New("invalid user id received from user service")
	}

	newAuth := &domain.AuthAccount{
		UserID:            userID,
		Email:             input.Email,
		Password:          hashedPwd,
		Provider:          ProviderLocal,
		ProviderAccountID: input.Email,
		IsEmailVerified:   false,
	}

	if err := u.authRepo.Register(ctx, newAuth); err != nil {
		return nil, fmt.Errorf("failed to register auth account: %w", err)
	}

	return &domain.AuthOutput{
		UserID:       userID.String(),
		AccessToken:  "",
		RefreshToken: "",
	}, nil
}

func (u *AuthUsecase) Login(ctx context.Context, input *domain.LoginInput) (*domain.AuthOutput, error) {
	existingUser, err := u.authRepo.FindByEmailAndProvider(ctx, input.Email, ProviderLocal)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return u.GenerateTokensAndSession(ctx, existingUser.UserID, existingUser.Email)
}

func (u *AuthUsecase) GetAllProviderAccounts(ctx context.Context, userID string) ([]domain.AuthAccount, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}

	accounts, err := u.authRepo.FindAllByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider accounts: %w", err)
	}

	return accounts, nil
}
