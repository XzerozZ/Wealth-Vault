package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"wealth-vault/user-service/internal/domain"
)

type MockMsgRepository struct {
	mock.Mock
}

func (m *MockMsgRepository) CreateMessage(
	ctx context.Context,
	log []domain.GroupMessage,
) error {

	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockMsgRepository) CreatePrivateMessage(
	ctx context.Context,
	log []domain.PrivateMessage,
) error {

	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockMsgRepository) GetGroupMessages(
	ctx context.Context,
	groupID string,
	userID string,
) ([]domain.GroupMessage, error) {

	args := m.Called(ctx, groupID, userID)

	var result []domain.GroupMessage
	if args.Get(0) != nil {
		result = args.Get(0).([]domain.GroupMessage)
	}

	return result, args.Error(1)
}

func (m *MockMsgRepository) GetPrivateMessages(
	ctx context.Context,
	userID, friendID string,
) ([]domain.PrivateMessage, error) {

	args := m.Called(ctx, userID, friendID)

	var result []domain.PrivateMessage
	if args.Get(0) != nil {
		result = args.Get(0).([]domain.PrivateMessage)
	}

	return result, args.Error(1)
}

func (m *MockMsgRepository) UpdateGrantMessageStatus(ctx context.Context, groupID, ownerID, targetID uuid.UUID, newMetadata string) error {
	args := m.Called(ctx, groupID, ownerID, targetID, newMetadata)
	return args.Error(0)
}

func (m *MockMsgRepository) CloseAllGrantPromptsForTarget(ctx context.Context, groupID, targetID uuid.UUID) error {
	args := m.Called(ctx, groupID, targetID)
	return args.Error(0)
}

func (m *MockMsgRepository) MarkAssetMessageAsDeletedinAssetService(ctx context.Context, assetID uuid.UUID) error {
	args := m.Called(ctx, assetID)
	return args.Error(0)
}

func (m *MockMsgRepository) MarkAllMemberAssetsAsUnshared(ctx context.Context, groupID, userID uuid.UUID) error {
	args := m.Called(ctx, groupID, userID)
	return args.Error(0)
}
