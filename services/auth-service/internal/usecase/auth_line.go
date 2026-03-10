package usecase

import (
	"context"
	"errors"
	"fmt"
	"wealth-vault/auth-service/internal/domain"

	"github.com/google/uuid"
)

func (u *AuthUsecase) LinkLineAccount(ctx context.Context, userID string, lineUserID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id format")
	}

	existing, _ := u.authRepo.FindByUserIDAndProvider(ctx, uid, ProviderLine)
	if existing != nil {
		return errors.New("user already linked with LINE account")
	}

	dummyEmail := fmt.Sprintf("line_%s@dummy.local", lineUserID)

	newAuth := &domain.AuthAccount{
		UserID:            uid,
		Provider:          ProviderLine,
		ProviderAccountID: lineUserID,
		Email:             dummyEmail,
		IsEmailVerified:   true,
	}

	if err := u.authRepo.Register(ctx, newAuth); err != nil {
		return fmt.Errorf("failed to link line account: %w", err)
	}

	return nil
}
