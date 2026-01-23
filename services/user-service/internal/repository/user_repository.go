package repository

import (
	"context"
	"wealth-vault/user-service/internal/domain"

	"github.com/google/uuid"
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

func (r *UserRepository) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
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

func (r *UserRepository) GetFriendList(ctx context.Context, userID uuid.UUID) ([]domain.FriendList, error) {
	var friendLists []domain.FriendList

	err := r.db.WithContext(ctx).Preload("Friend").Where("user_id = ?", userID).Find(&friendLists).Error
	if err != nil {
		return nil, err
	}

	return friendLists, nil
}

func (r *UserRepository) AddFriend(ctx context.Context, fri *domain.FriendList) error {
	if err := r.db.Create(&fri).Error; err != nil {
		return err
	}

	return nil
}
