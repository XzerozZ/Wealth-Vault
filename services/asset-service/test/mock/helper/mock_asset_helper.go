package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"wealth-vault/asset-service/internal/domain"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Delete(url string) error {
	args := m.Called(url)
	return args.Error(0)
}

type MockAssetHelper struct {
	mock.Mock
}

func (m *MockAssetHelper) SyncFiles(
	ctx context.Context,
	params domain.FileSyncParams,
) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *MockAssetHelper) CleanupResource(
	ctx context.Context,
	entityID uuid.UUID,
	files []domain.FileAssociate,
	hardDeleteFunc func(uuid.UUID) error,
) {
	m.Called(ctx, entityID, files, hardDeleteFunc)
}
