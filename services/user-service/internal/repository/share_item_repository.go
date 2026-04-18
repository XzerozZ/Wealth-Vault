package repository

import (
	"context"
	"fmt"
	"strings"
	"time"
	"wealth-vault/user-service/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
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
	if err := r.db.WithContext(ctx).Where("id = ? AND owner_id = ?", itemID, userID).Delete(&domain.GroupItem{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *ShareItemRepository) DeleteIteminFriend(ctx context.Context, itemID uuid.UUID, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND owner_id = ?", itemID, userID).Delete(&domain.FriendItem{}).Error; err != nil {
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

func (r *ShareItemRepository) GetItemSharedTargets(ctx context.Context, userID, itemID uuid.UUID, itemType string) (*domain.SharedTargetsResult, error) {
	result := &domain.SharedTargetsResult{
		Groups:  make([]domain.SharedGroupTarget, 0),
		Friends: make([]domain.SharedFriendTarget, 0),
		Emails:  make([]domain.SharedEmailTarget, 0),
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		query := `
			groups.id as group_id, 
			groups.group_name as group_name, 
			groups.group_profile as group_image, 
			group_items.share_at as shared_at,
			(SELECT COUNT(*) FROM group_members WHERE group_members.group_id = groups.id) as member_count
		`

		return r.db.WithContext(ctx).
			Table("group_items").
			Select(query).
			Joins("JOIN groups ON groups.id = group_items.group_id").
			Where("group_items.owner_id = ? AND group_items.entity_id = ? AND group_items.entity_type = ?", userID, itemID, itemType).
			Scan(&result.Groups).Error
	})

	g.Go(func() error {
		return r.db.WithContext(ctx).
			Table("friend_items").
			Select("users.id as friend_id, users.username, users.profile, friend_items.share_at as shared_at").
			Joins("JOIN users ON users.id = friend_items.friend_id").
			Where("friend_items.owner_id = ? AND friend_items.entity_id = ? AND friend_items.entity_type = ?", userID, itemID, itemType).
			Scan(&result.Friends).Error
	})

	g.Go(func() error {
		return r.db.WithContext(ctx).
			Table("email_items").
			Select("DISTINCT ON (email) email, share_at as shared_at, is_sent").
			Where("owner_id = ? AND entity_id = ? AND entity_type = ?", userID, itemID, itemType).
			Order("email, share_at DESC").
			Scan(&result.Emails).Error
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ShareItemRepository) GetSharedItemIDs(ctx context.Context, userID, targetID uuid.UUID, targetType string) ([]string, error) {
	var ids []string
	var err error

	normalizedType := strings.ToUpper(targetType)
	if normalizedType == "GROUP" {
		if err = r.db.WithContext(ctx).Table("group_items").
			Where("owner_id = ? AND group_id = ?", userID, targetID).
			Pluck("entity_id", &ids).Error; err != nil {
			return nil, err
		}
	} else if normalizedType == "FRIEND" {
		if err = r.db.WithContext(ctx).
			Table("friend_items").
			Where("owner_id = ? AND friend_id = ?", userID, targetID).
			Pluck("entity_id", &ids).Error; err != nil {
			return nil, err
		}
	}

	return ids, nil
}

func (r *ShareItemRepository) GetItemsSharedByFriend(ctx context.Context, myUserID, friendID uuid.UUID) ([]domain.SharedItemSummary, error) {
	var items []domain.SharedItemSummary

	query := `
		SELECT entity_id, entity_type 
		FROM friend_items 
		WHERE owner_id = ? AND friend_id = ? 

		UNION

		SELECT gi.entity_id, gi.entity_type 
		FROM group_items gi
		INNER JOIN group_item_viewers giv ON gi.id = giv.group_item_id
		WHERE gi.owner_id = ? 
		AND giv.viewer_id = ?
	`

	err := r.db.WithContext(ctx).Raw(query, friendID, myUserID, friendID, myUserID).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ShareItemRepository) GetAllSharedItemIDsByUser(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var itemIDs []string
	query := `
		SELECT entity_id FROM group_items WHERE owner_id = ?
		UNION
		SELECT entity_id FROM friend_items WHERE owner_id = ?
		UNION
		SELECT entity_id FROM email_items WHERE owner_id = ?
	`
	if err := r.db.WithContext(ctx).Raw(query, userID, userID, userID).Pluck("entity_id", &itemIDs).Error; err != nil {
		return nil, err
	}

	return itemIDs, nil
}

func (r *ShareItemRepository) GetGroupItemByID(ctx context.Context, id uuid.UUID) (*domain.GroupItem, error) {
	var item domain.GroupItem
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ShareItemRepository) GetFriendItemByID(ctx context.Context, id uuid.UUID) (*domain.FriendItem, error) {
	var item domain.FriendItem
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
