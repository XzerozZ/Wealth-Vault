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
