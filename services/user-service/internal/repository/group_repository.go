package repository

import (
	"context"
	"wealth-vault/user-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) CreateGroup(ctx context.Context, group *domain.Group, memberIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(group).Error; err != nil {
			return err
		}

		members := []domain.GroupMember{}

		members = append(members, domain.GroupMember{
			GroupID: group.ID,
			UserID:  group.CreatedBy,
			Role:    "admin",
		})

		for _, memberIDStr := range memberIDs {
			mid, err := uuid.Parse(memberIDStr)
			if err != nil {
				continue
			}

			if mid == group.CreatedBy {
				continue
			}

			members = append(members, domain.GroupMember{
				GroupID: group.ID,
				UserID:  mid,
				Role:    "member",
			})
		}

		if len(members) > 0 {
			if err := tx.Create(&members).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *GroupRepository) GetMember(ctx context.Context, id uuid.UUID) ([]*domain.User, int64, error) {
	var mem []*domain.User
	var total int64
	query := r.db.Model(&domain.User{}).
		Joins("JOIN group_members ON group_members.user_id = users.id").
		Where("group_members.group_id = ?", id)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("users.username ASC").Find(&mem).Error; err != nil {
		return nil, 0, err
	}

	return mem, total, nil
}

func (r *GroupRepository) GetGroup(ctx context.Context, id uuid.UUID) (*domain.Group, int64, error) {
	var group domain.Group
	var total int64
	if err := r.db.WithContext(ctx).Model(&domain.GroupMember{}).Where("group_id = ?", id).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&group).Error; err != nil {
		return nil, 0, err
	}
	return &group, total, nil
}

func (r *GroupRepository) AllGetGroup(ctx context.Context, uid uuid.UUID) ([]domain.GroupWithCount, error) {
	var results []domain.GroupWithCount
	err := r.db.WithContext(ctx).
		Table("groups g").
		Select(`
			g.*,
			(SELECT COUNT(*) FROM group_members WHERE group_id = g.id) as member_count
		`).
		Joins("JOIN group_members gm ON gm.group_id = g.id").
		Where("gm.user_id = ?", uid).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *GroupRepository) IsUserMember(ctx context.Context, id uuid.UUID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.GroupMember{}).
		Where("group_id = ? AND user_id = ?", id, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *GroupRepository) UpdateGroup(ctx context.Context, group *domain.Group, mask []string, logEntry *domain.GroupLog) (*domain.Group, int64, error) {
	var total int64
	query := r.db.Model(&domain.User{}).
		Joins("JOIN group_members ON group_members.user_id = users.id").
		Where("group_members.group_id = ?", group.ID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	tx := r.db.WithContext(ctx).Model(group).Where("id = ?", group.ID)
	if len(mask) > 0 {
		tx = tx.Select(mask)
	}

	if err := tx.Updates(group).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Create(logEntry).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.First(group, "id = ?", group.ID).Error; err != nil {
		return nil, 0, err
	}

	return group, total, nil
}

func (r *GroupRepository) RemoveMemberAndTheirSharedItems(ctx context.Context, groupID, memberID uuid.UUID, logEntry *domain.GroupLog) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ? AND owner_id = ?", groupID, memberID).Delete(&domain.GroupItem{}).Error; err != nil {
			return err
		}

		if err := tx.Where("group_id = ? AND user_id = ?", groupID, memberID).Delete(&domain.GroupMember{}).Error; err != nil {
			return err
		}

		if err := tx.Create(logEntry).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *GroupRepository) CreateLog(ctx context.Context, log *domain.GroupLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *GroupRepository) DeleteGroupCompletely(ctx context.Context, groupID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ?", groupID).Delete(&domain.GroupItem{}).Error; err != nil {
			return err
		}

		if err := tx.Where("group_id = ?", groupID).Delete(&domain.GroupMember{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", groupID).Delete(&domain.Group{}).Error; err != nil {
			return err
		}

		return nil
	})
}
