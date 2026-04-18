package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"wealth-vault/user-service/internal/domain"
)

type MockShareItemRepository struct {
	mock.Mock
}

func (m *MockShareItemRepository) ShareItemtoGroup(ctx context.Context, items []domain.GroupItem) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *MockShareItemRepository) ShareItemtoFriend(ctx context.Context, items []domain.FriendItem) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *MockShareItemRepository) ShareItemtoEmail(ctx context.Context, items []domain.EmailItem) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}

func (m *MockShareItemRepository) GetExistingSharedMap(
	ctx context.Context,
	ownerID, friendID uuid.UUID,
) (map[string]bool, error) {

	args := m.Called(ctx, ownerID, friendID)

	var result map[string]bool
	if args.Get(0) != nil {
		result = args.Get(0).(map[string]bool)
	}

	return result, args.Error(1)
}

func (m *MockShareItemRepository) IsItemSharedtoGroup(
	ctx context.Context,
	groupID, entityID uuid.UUID,
	entityType string,
) (bool, error) {

	args := m.Called(ctx, groupID, entityID, entityType)
	return args.Bool(0), args.Error(1)
}

func (m *MockShareItemRepository) IsItemSharedtoFriend(
	ctx context.Context,
	friendID, entityID uuid.UUID,
	entityType string,
) (bool, error) {

	args := m.Called(ctx, friendID, entityID, entityType)
	return args.Bool(0), args.Error(1)
}

func (m *MockShareItemRepository) IsItemSharedtoEmail(
	ctx context.Context,
	entityID uuid.UUID,
	email, entityType string,
) (bool, error) {

	args := m.Called(ctx, entityID, email, entityType)
	return args.Bool(0), args.Error(1)
}

func (m *MockShareItemRepository) IsGroupMember(
	ctx context.Context,
	groupID uuid.UUID,
	userID uuid.UUID,
) (bool, error) {

	args := m.Called(ctx, groupID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockShareItemRepository) GetSharedIteminGroup(
	ctx context.Context,
	groupID, userID uuid.UUID,
) ([]domain.GroupItem, error) {

	args := m.Called(ctx, groupID, userID)

	var result []domain.GroupItem
	if args.Get(0) != nil {
		result = args.Get(0).([]domain.GroupItem)
	}

	return result, args.Error(1)
}

func (m *MockShareItemRepository) GetSharedIteminFriend(
	ctx context.Context,
	friendID, userID uuid.UUID,
) ([]domain.FriendItem, error) {

	args := m.Called(ctx, friendID, userID)

	var result []domain.FriendItem
	if args.Get(0) != nil {
		result = args.Get(0).([]domain.FriendItem)
	}

	return result, args.Error(1)
}

func (m *MockShareItemRepository) GetOwnedItemIDs(
	ctx context.Context,
	itemIDs []string,
	ownerID uuid.UUID,
) ([]uuid.UUID, error) {

	args := m.Called(ctx, itemIDs, ownerID)

	var result []uuid.UUID
	if args.Get(0) != nil {
		result = args.Get(0).([]uuid.UUID)
	}

	return result, args.Error(1)
}

func (m *MockShareItemRepository) GetFutureItemsInGroup(
	ctx context.Context,
	groupID uuid.UUID,
) ([]uuid.UUID, error) {

	args := m.Called(ctx, groupID)

	var result []uuid.UUID
	if args.Get(0) != nil {
		result = args.Get(0).([]uuid.UUID)
	}

	return result, args.Error(1)
}

func (m *MockShareItemRepository) GetItemOwnersInGroup(
	ctx context.Context,
	groupID uuid.UUID,
) ([]uuid.UUID, error) {

	args := m.Called(ctx, groupID)

	var result []uuid.UUID
	if args.Get(0) != nil {
		result = args.Get(0).([]uuid.UUID)
	}

	return result, args.Error(1)
}

func (m *MockShareItemRepository) GetItemSharedTargets(
	ctx context.Context,
	userID, itemID uuid.UUID,
	itemType string,
) (*domain.SharedTargetsResult, error) {

	args := m.Called(ctx, userID, itemID, itemType)

	var result *domain.SharedTargetsResult
	if args.Get(0) != nil {
		result = args.Get(0).(*domain.SharedTargetsResult)
	}

	return result, args.Error(1)
}

func (m *MockShareItemRepository) GetSharedItemIDs(
	ctx context.Context,
	userID, targetID uuid.UUID,
	targetType string,
) ([]string, error) {

	args := m.Called(ctx, userID, targetID, targetType)

	var result []string
	if args.Get(0) != nil {
		result = args.Get(0).([]string)
	}

	return result, args.Error(1)
}

func (m *MockShareItemRepository) DeleteIteminGroup(
	ctx context.Context,
	itemID uuid.UUID,
	userID uuid.UUID,
) error {

	args := m.Called(ctx, itemID, userID)
	return args.Error(0)
}

func (m *MockShareItemRepository) DeleteIteminFriend(
	ctx context.Context,
	itemID uuid.UUID,
	userID uuid.UUID,
) error {

	args := m.Called(ctx, itemID, userID)
	return args.Error(0)
}

func (m *MockShareItemRepository) GetPendingEmails(
	ctx context.Context,
) ([]domain.EmailItem, error) {

	args := m.Called(ctx)

	var result []domain.EmailItem
	if args.Get(0) != nil {
		result = args.Get(0).([]domain.EmailItem)
	}

	return result, args.Error(1)
}

func (m *MockShareItemRepository) MarkEmailsAsSent(
	ctx context.Context,
	ids []uuid.UUID,
) error {

	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockShareItemRepository) CountItemsByOwner(
	ctx context.Context,
	itemIDs []string,
	ownerID uuid.UUID,
) (int64, error) {

	args := m.Called(ctx, itemIDs, ownerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockShareItemRepository) AddMember(
	ctx context.Context,
	members []domain.GroupMember,
) error {

	args := m.Called(ctx, members)
	return args.Error(0)
}

func (m *MockShareItemRepository) BatchCreateViewers(
	ctx context.Context,
	viewers []domain.GroupItemViewer,
) error {

	args := m.Called(ctx, viewers)
	return args.Error(0)
}

func (m *MockShareItemRepository) DeleteAllReferencesByEntityID(
	ctx context.Context,
	entityID uuid.UUID,
) error {

	args := m.Called(ctx, entityID)
	return args.Error(0)
}

func (m *MockShareItemRepository) GetItemsSharedByFriend(
	ctx context.Context,
	myUserID, friendID uuid.UUID,
) ([]domain.SharedItemSummary, error) {

	args := m.Called(ctx, myUserID, friendID)

	var result []domain.SharedItemSummary
	if args.Get(0) != nil {
		result = args.Get(0).([]domain.SharedItemSummary)
	}

	return result, args.Error(1)
}

func (m *MockShareItemRepository) GetAllSharedItemIDsByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockShareItemRepository) GetFriendItemByID(ctx context.Context, id uuid.UUID) (*domain.FriendItem, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.FriendItem), args.Error(1)
}
func (m *MockShareItemRepository) GetGroupItemByID(ctx context.Context, id uuid.UUID) (*domain.GroupItem, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.GroupItem), args.Error(1)
}
