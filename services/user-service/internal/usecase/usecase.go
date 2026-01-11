package usecase

import (
	"context"
	"time"
	"wealth-vault/user-service/internal/domain"
	repo "wealth-vault/user-service/internal/repository/interface"

	"github.com/google/uuid"
)

type UserUsecase struct {
	userRepo repo.UserRepository
}

func NewUserUsecase(r repo.UserRepository) UserUsecase {
	return UserUsecase{userRepo: r}
}

func (u *UserUsecase) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	user.ID = uuid.NewString()
	if err := u.userRepo.CreateUser(ctx, user); err != nil {
		return "", err
	}

	return user.ID, nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, input *domain.UpdateUserInput) (*domain.User, error) {
	updateData := &domain.User{
		ID:          input.ID,
		Firstname:   input.Firstname,
		Lastname:    input.Lastname,
		Username:    input.Username,
		Profile:     input.Profile,
		Phonenumber: input.Phonenumber,
	}

	if input.BirthdayStr != "" {
		parsedTime, err := time.Parse("2006-01-02", input.BirthdayStr)
		if err == nil {
			updateData.Birthday = parsedTime
		}
	}

	updatedUser, err := u.userRepo.UpdateUser(ctx, updateData, input.UpdateMask)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
