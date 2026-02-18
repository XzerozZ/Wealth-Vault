package mock

import (
	"context"
	"wealth-vault/asset-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) CreateFiles(ctx context.Context, files []domain.FileAssociate) error {
	args := m.Called(ctx, files)
	return args.Error(0)
}

func (m *MockFileRepository) DeleteFiles(ctx context.Context, fileIDs []uuid.UUID) error {
	args := m.Called(ctx, fileIDs)
	return args.Error(0)
}

func (m *MockFileRepository) GetFilesByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.FileAssociate, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]domain.FileAssociate), args.Error(1)
}
