package repository

import (
	"context"
	"wealth-vault/user-service/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User, mask []string) (*domain.User, error) {
	tx := r.db.WithContext(ctx).Model(user).Where("id = ?", user.ID)
	if len(mask) > 0 {
		tx = tx.Select(mask)
	}

	if err := tx.Updates(user).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(user, "id = ?", user.ID).Error; err != nil {
		return nil, err
	}

	return user, nil
}
