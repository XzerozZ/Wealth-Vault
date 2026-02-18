package mapper_test

import (
	"testing"

	"wealth-vault/asset-service/internal/domain"
	"wealth-vault/asset-service/pkg/utils/mapper"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToPbFiles(t *testing.T) {
	fileID := uuid.New()

	tests := []struct {
		name  string
		input []domain.FileAssociate
	}{
		{
			name: "Success - Map Files to Proto",
			input: []domain.FileAssociate{
				{ID: fileID, Link: "https://vault.com/s/1.png", FileType: "image/png"},
			},
		},
		{
			name:  "Edge Case - Empty Input",
			input: []domain.FileAssociate{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToPbFiles(tt.input)

			if len(tt.input) == 0 {
				assert.NotNil(t, result)
				assert.Empty(t, result)
			} else {
				assert.Len(t, result, 1)
				assert.Equal(t, fileID.String(), result[0].Id)
				assert.Equal(t, tt.input[0].Link, result[0].Url)
			}
		})
	}
}

func TestToReferenceProtos(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()

	t.Run("ToRefLandProto - Valid Mapping", func(t *testing.T) {
		input := []domain.Land{
			{ID: id1, Name: "Land A"},
		}
		result := mapper.ToRefLandProto(input)
		assert.Len(t, result, 1)
		assert.Equal(t, id1.String(), result[0].Id)
		assert.Equal(t, "Land A", result[0].Name)
	})

	t.Run("ToRefInsProto - Valid Mapping", func(t *testing.T) {
		input := []domain.Insurance{
			{ID: id2, Name: "Life Ins"},
		}
		result := mapper.ToRefInsProto(input)
		assert.Len(t, result, 1)
		assert.Equal(t, id2.String(), result[0].Id)
	})

	t.Run("Empty Input - Should Return Empty Slice Not Nil", func(t *testing.T) {
		result := mapper.ToRefBuildingProto([]domain.Building{})
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})
}
