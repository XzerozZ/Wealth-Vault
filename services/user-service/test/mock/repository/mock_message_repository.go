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
) ([]domain.GroupMessage, error) {

	args := m.Called(ctx, groupID)

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

func (m *MockMsgRepository) UpdateGrantMessageStatus(ctx context.Context, groupID, targetID uuid.UUID, newMetadata string) error {
	args := m.Called(ctx, groupID, targetID, newMetadata)
	return args.Error(0)
}
