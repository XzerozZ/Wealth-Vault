// File: auth_oauth.go
package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"wealth-vault/auth-service/internal/domain"
	pb "wealth-vault/auth-service/pkg/pb/proto/user"
)

func (u *AuthUsecase) LoginWithGoogle(ctx context.Context, googleIDToken string) (*domain.AuthOutput, error) {
	payload, err := u.googleValidator.Validate(ctx, googleIDToken)
	if err != nil {
		return nil, fmt.Errorf("invalid google id token: %w", err)
	}

	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return nil, errors.New("email not found in google token")
	}

	existingUser, err := u.authRepo.FindByEmailAndProvider(ctx, email, ProviderGoogle)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	var userID uuid.UUID

	if existingUser != nil {
		userID = existingUser.UserID
	} else {
		name, _ := payload.Claims["name"].(string)
		if name == "" {
			name = strings.Split(email, "@")[0]
		}

		userRes, err := u.userClient.CreateUser(ctx, &pb.CreateUserRequest{
			Email:    email,
			Username: name,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create user profile: %w", err)
		}

		userID, err = uuid.Parse(userRes.User.Id)
		if err != nil {
			return nil, errors.New("invalid user id from user service")
		}

		newAuth := &domain.AuthAccount{
			UserID:            userID,
			Email:             email,
			Password:          "",
			Provider:          ProviderGoogle,
			ProviderAccountID: payload.Subject,
			IsEmailVerified:   true,
		}

		if err := u.authRepo.Register(ctx, newAuth); err != nil {
			return nil, fmt.Errorf("failed to register auth account: %w", err)
		}
	}

	return u.GenerateTokensAndSession(ctx, userID, email)
}
