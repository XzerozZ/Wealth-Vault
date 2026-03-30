package repository

import (
	"context"
	"wealth-vault/user-service/internal/domain"

	"github.com/google/uuid"
)

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

func (r *UserRepository) GetCloseFriends(ctx context.Context, userID uuid.UUID) ([]domain.FriendList, error) {
	var friendLists []domain.FriendList
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_close_friend = ?", userID, true).
		Preload("Friend").
		Find(&friendLists).Error

	if err != nil {
		return nil, err
	}

	return friendLists, nil
}
