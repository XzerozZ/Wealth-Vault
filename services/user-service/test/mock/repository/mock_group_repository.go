package mock

import (
	"context"
	"wealth-vault/user-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockGroupRepository struct {
	mock.Mock
}

func (m *MockGroupRepository) CreateGroup(
	ctx context.Context,
	group *domain.Group,
	initialMembers []string,
) error {

	args := m.Called(ctx, group, initialMembers)
	return args.Error(0)
}

func (m *MockGroupRepository) CreateLog(
	ctx context.Context,
	log *domain.GroupLog,
) error {

	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockGroupRepository) GetMember(
	ctx context.Context,
	id uuid.UUID,
) ([]*domain.User, int64, error) {

	args := m.Called(ctx, id)

	var users []*domain.User
	if args.Get(0) != nil {
		users = args.Get(0).([]*domain.User)
	}

	return users, args.Get(1).(int64), args.Error(2)
}

func (m *MockGroupRepository) GetGroup(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Group, int64, error) {

	args := m.Called(ctx, id)

	var group *domain.Group
	if args.Get(0) != nil {
		group = args.Get(0).(*domain.Group)
	}

	return group, args.Get(1).(int64), args.Error(2)
}

func (m *MockGroupRepository) AllGetGroup(
	ctx context.Context,
	uid uuid.UUID,
) ([]domain.GroupWithCount, error) {

	args := m.Called(ctx, uid)

	var groups []domain.GroupWithCount
	if args.Get(0) != nil {
		groups = args.Get(0).([]domain.GroupWithCount)
	}

	return groups, args.Error(1)
}

func (m *MockGroupRepository) IsUserMember(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
) (bool, error) {

	args := m.Called(ctx, id, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockGroupRepository) UpdateGroup(
	ctx context.Context,
	group *domain.Group,
	mask []string,
	logEntry *domain.GroupLog,
) (*domain.Group, int64, error) {

	args := m.Called(ctx, group, mask, logEntry)

	var updated *domain.Group
	if args.Get(0) != nil {
		updated = args.Get(0).(*domain.Group)
	}

	return updated, args.Get(1).(int64), args.Error(2)
}

func (m *MockGroupRepository) RemoveMemberAndTheirSharedItems(
	ctx context.Context,
	groupID, memberID uuid.UUID,
	logEntry *domain.GroupLog,
) error {

	args := m.Called(ctx, groupID, memberID, logEntry)
	return args.Error(0)
}

func (m *MockGroupRepository) DeleteGroupCompletely(
	ctx context.Context,
	groupID uuid.UUID,
) error {

	args := m.Called(ctx, groupID)
	return args.Error(0)
}
