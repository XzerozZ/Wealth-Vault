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

func (r *UserRepository) GetFriendList(ctx context.Context, userID uuid.UUID) ([]domain.FriendList, error) {
	var friendLists []domain.FriendList
	err := r.db.WithContext(ctx).Where("user_id = ? AND status = ?", userID, "ACCEPTED").Preload("Friend").Find(&friendLists).Error
	if err != nil {
		return nil, err
	}

	return friendLists, nil
}

func (r *UserRepository) AddFriend(ctx context.Context, fri *domain.FriendList) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.FriendList{}).
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
			fri.UserID, fri.FriendID, fri.FriendID, fri.UserID).Count(&count).Error; err != nil {
		return err
	}

	if err := r.db.Create(&fri).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) CreateFriendship(ctx context.Context, fri *domain.FriendList) error {
	if err := r.db.WithContext(ctx).Create(&fri).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
			userID, friendID, friendID, userID).
		Delete(&domain.FriendList{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateFriendStatus(ctx context.Context, userID, friendID uuid.UUID, status string) error {
	if err := r.db.WithContext(ctx).Model(&domain.FriendList{}).Where("user_id = ? AND friend_id = ?", friendID, userID).Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) CheckFriendship(ctx context.Context, userID, friendID uuid.UUID) (bool, string, error) {
	var friendList domain.FriendList
	if err := r.db.WithContext(ctx).Where("user_id = ? AND friend_id = ?", userID, friendID).First(&friendList).Error; err != nil {
		return false, "", err
	}

	return true, friendList.Status, nil
}

func (r *UserRepository) GetIncomingRequests(ctx context.Context, userID uuid.UUID) ([]domain.FriendList, error) {
	var requests []domain.FriendList
	if err := r.db.WithContext(ctx).Where("friend_id = ? AND status = ?", userID, "PENDING").Preload("User").Find(&requests).Error; err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *UserRepository) SetCloseFriendStatus(ctx context.Context, userID, friendID uuid.UUID, isClose bool) error {
	if err := r.db.WithContext(ctx).Model(&domain.FriendList{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Update("is_close_friend", isClose).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetCloseFriends(ctx context.Context, userID uuid.UUID) ([]domain.User, error) {
	var friends []domain.User
	err := r.db.WithContext(ctx).
		Model(&domain.User{ID: userID}).
		Association("Friends").
		Find(&friends, "friend_lists.is_close_friend = ?", true)

	if err != nil {
		return nil, err
	}

	return friends, nil
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
