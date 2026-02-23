package mock

import (
	"github.com/stretchr/testify/mock"
)

type MockSupabaseStorage struct {
	mock.Mock
}

func (m *MockSupabaseStorage) Delete(fileURL string) error {
	args := m.Called(fileURL)
	return args.Error(0)
}
