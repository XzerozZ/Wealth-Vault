package mapper_test

import (
	"testing"

	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils/mapper"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToDomainFiles(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name     string
		input    []*pb.FileInfo
		expected []domain.FileAssociate
	}{
		{
			name: "Success - Map Multiple Files",
			input: []*pb.FileInfo{
				{Url: "https://storage.com/file1.jpg", FileType: "image/jpeg"},
				{Url: "https://storage.com/file2.pdf", FileType: "application/pdf"},
			},
			expected: []domain.FileAssociate{
				{Link: "https://storage.com/file1.jpg", FileType: "image/jpeg", UserID: userID},
				{Link: "https://storage.com/file2.pdf", FileType: "application/pdf", UserID: userID},
			},
		},
		{
			name:     "Edge Case - Empty Slice",
			input:    []*pb.FileInfo{},
			expected: nil,
		},
		{
			name:     "Edge Case - Nil Input",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToDomainFiles(tt.input, userID)

			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.Len(t, result, len(tt.expected))
				for i := range result {
					assert.Equal(t, tt.expected[i].Link, result[i].Link)
					assert.Equal(t, tt.expected[i].FileType, result[i].FileType)
					assert.Equal(t, tt.expected[i].UserID, result[i].UserID)
				}
			}
		})
	}
}
