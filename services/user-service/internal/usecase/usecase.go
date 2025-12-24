package usecase

import (
	"context"
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

func (u *UserUsecase) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return u.userRepo.FindByID(ctx, id)
}
