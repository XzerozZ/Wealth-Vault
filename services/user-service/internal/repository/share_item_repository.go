package repository

import (
	"context"
	"fmt"
	"time"
	"wealth-vault/user-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShareItemRepository struct {
	db *gorm.DB
}

func NewShareItemRepository(db *gorm.DB) *ShareItemRepository {
	return &ShareItemRepository{db: db}
}

func (r *ShareItemRepository) ShareItemtoGroup(ctx context.Context, items []domain.GroupItem) error {
	if len(items) == 0 {
		return nil
	}

	if err := r.db.Create(&items).Error; err != nil {
		return err
	}

	return nil
}

func (r *ShareItemRepository) ShareItemtoFriend(ctx context.Context, items []domain.FriendItem) error {
	if len(items) == 0 {
		return nil
	}

	if err := r.db.Create(&items).Error; err != nil {
		return err
	}

	return nil
}

func (r *ShareItemRepository) ShareItemtoEmail(ctx context.Context, items []domain.EmailItem) error {
	if len(items) == 0 {
		return nil
	}

	if err := r.db.Create(&items).Error; err != nil {
		return err
	}

	return nil
}

func (r *ShareItemRepository) GetExistingSharedMap(ctx context.Context, ownerID, friendID uuid.UUID) (map[string]bool, error) {
	type Result struct {
		EntityID   uuid.UUID
		EntityType string
	}
	var results []Result
	err := r.db.WithContext(ctx).Table("friend_items").
		Select("entity_id, entity_type").
		Where("owner_id = ? AND friend_id = ?", ownerID, friendID).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	m := make(map[string]bool)
	for _, res := range results {
		key := fmt.Sprintf("%s:%s", res.EntityType, res.EntityID.String())
		m[key] = true
	}

	return m, nil
}

func (r *ShareItemRepository) IsItemSharedtoGroup(ctx context.Context, groupID, entityID uuid.UUID, entityType string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.GroupItem{}).Where("group_id = ? AND entity_type = ? AND entity_id = ?", groupID, entityType, entityID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ShareItemRepository) IsItemSharedtoFriend(ctx context.Context, friendID, entityID uuid.UUID, entityType string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.FriendItem{}).Where("friend_id = ? AND entity_type = ? AND entity_id = ?", friendID, entityType, entityID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ShareItemRepository) IsItemSharedtoEmail(ctx context.Context, entityID uuid.UUID, email, entityType string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.EmailItem{}).Where("email = ? AND entity_type = ? AND entity_id = ?", email, entityType, entityID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ShareItemRepository) GetSharedIteminGroup(ctx context.Context, groupID, userID uuid.UUID) ([]domain.GroupItem, error) {
	var items []domain.GroupItem
	err := r.db.WithContext(ctx).
		Joins("LEFT JOIN group_item_viewers v ON v.group_item_id = group_items.id").
		Where("group_items.group_id = ?", groupID).
		Where("group_items.owner_id = ? OR v.viewer_id = ?", userID, userID).
		Where("group_items.share_at <= ?", time.Now()).
		Order("group_items.share_at DESC").
		Preload("User").
		Distinct().
		Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ShareItemRepository) GetSharedIteminFriend(ctx context.Context, friendID, userID uuid.UUID) ([]domain.FriendItem, error) {
	var items []domain.FriendItem
	err := r.db.WithContext(ctx).Where(
		r.db.Where("owner_id = ? AND friend_id = ?", userID, friendID).Or("owner_id = ? AND friend_id = ?", friendID, userID),
	).Where("share_at <= ?", time.Now()).Order("share_at DESC").Preload("User").Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ShareItemRepository) DeleteIteminGroup(ctx context.Context, itemID uuid.UUID, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND created_by = ?", itemID, userID).Delete(&domain.GroupItem{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *ShareItemRepository) DeleteIteminFriend(ctx context.Context, itemID uuid.UUID, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND created_by = ?", itemID, userID).Delete(&domain.FriendItem{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *ShareItemRepository) GetPendingEmails(ctx context.Context) ([]domain.EmailItem, error) {
	var items []domain.EmailItem
	if err := r.db.WithContext(ctx).Where("share_at <= ? AND is_sent = ?", time.Now(), false).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ShareItemRepository) MarkEmailsAsSent(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).Model(&domain.EmailItem{}).Where("id IN ?", ids).Update("is_sent", true).Error; err != nil {
		return err
	}

	return nil
}

func (r *ShareItemRepository) IsGroupMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ShareItemRepository) CountItemsByOwner(ctx context.Context, itemIDs []string, ownerID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.GroupItem{}).Where("id IN ? AND owner_id = ?", itemIDs, ownerID).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ShareItemRepository) GetOwnedItemIDs(ctx context.Context, itemIDs []string, ownerID uuid.UUID) ([]uuid.UUID, error) {
	var validIDs []uuid.UUID
	err := r.db.Table("group_items").
		Where("id IN ? AND owner_id = ?", itemIDs, ownerID).
		Pluck("id", &validIDs).
		Error

	return validIDs, err
}

func (r *ShareItemRepository) AddMember(ctx context.Context, members []domain.GroupMember) error {
	if len(members) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "group_id"}, {Name: "user_id"}},
			DoNothing: true,
		}).
		Create(&members).Error; err != nil {
		return err
	}

	return nil
}

func (r *ShareItemRepository) BatchCreateViewers(ctx context.Context, viewers []domain.GroupItemViewer) error {
	if len(viewers) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "group_item_id"}, {Name: "viewer_id"}},
			DoNothing: true,
		}).Create(&viewers).Error; err != nil {
		return err
	}

	return nil
}

func (r *ShareItemRepository) GetFutureItemsInGroup(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := r.db.WithContext(ctx).
		Model(&domain.GroupItem{}).
		Where("group_id = ? AND share_at > ?", groupID, time.Now()).
		Pluck("id", &ids).Error
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *ShareItemRepository) GetItemOwnersInGroup(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error) {
	var ownerIDs []uuid.UUID
	if err := r.db.WithContext(ctx).Model(&domain.GroupItem{}).
		Where("group_id = ?", groupID).Distinct("owner_id").
		Pluck("owner_id", &ownerIDs).Error; err != nil {
		return nil, err
	}

	return ownerIDs, nil
}

func (r *ShareItemRepository) DeleteAllReferencesByEntityID(ctx context.Context, entityID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("entity_id = ?", entityID).Delete(&domain.GroupItem{}).Error; err != nil {
			return err
		}

		if err := tx.Where("entity_id = ?", entityID).Delete(&domain.FriendItem{}).Error; err != nil {
			return err
		}

		if err := tx.Where("entity_id = ?", entityID).Delete(&domain.EmailItem{}).Error; err != nil {
			return err
		}

		return nil
	})
}
