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
	if err := r.db.Create(user).Error; err != nil {
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

func (r *UserRepository) GetUsersReadyForAutoShare(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	if err := r.db.WithContext(ctx).
		Where("is_auto_share_enabled = ? AND is_auto_share_triggered = ?", true, false).
		Where("EXTRACT(YEAR FROM age(birthday)) >= auto_share_age").
		Preload("Friends", func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN friend_lists ON friend_lists.friend_id = users.id").
				Where("friend_lists.status = ? AND friend_lists.is_close_friend = ?", "ACCEPTED", true)
		}).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) MarkAutoShareTriggered(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", userID).
		Update("is_auto_share_triggered", true).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) CreateFriendLog(ctx context.Context, log *domain.FriendLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
