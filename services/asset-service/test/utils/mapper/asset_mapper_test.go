package mapper_test

import (
	"testing"
	"time"

	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils/mapper"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToAssetSummaryProto(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	tests := []struct {
		name     string
		input    domain.AssetSummary
		expected *pb.AssetSummary
	}{
		{
			name: "Success - Full Data Mapping",
			input: domain.AssetSummary{
				ID:        id,
				Type:      "CASH",
				Name:      "Emergency Fund",
				Value:     50000.00,
				CreatedAt: now,
			},
			expected: &pb.AssetSummary{
				Id:    id.String(),
				Type:  "CASH",
				Name:  "Emergency Fund",
				Value: 50000.00,
			},
		},
		{
			name: "Success - Empty Values",
			input: domain.AssetSummary{
				ID: id,
			},
			expected: &pb.AssetSummary{
				Id: id.String(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToAssetSummaryProto(tt.input)

			assert.Equal(t, tt.expected.Id, result.Id)
			assert.Equal(t, tt.expected.Type, result.Type)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Value, result.Value)

			assert.NotNil(t, result.CreatedAt)
			assert.Equal(t, tt.input.CreatedAt.Unix(), result.CreatedAt.AsTime().Unix())
		})
	}
}

func TestToAssetSummaryProtoList(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()

	t.Run("Should map multiple items correctly", func(t *testing.T) {
		domains := []domain.AssetSummary{
			{ID: id1, Name: "Asset 1"},
			{ID: id2, Name: "Asset 2"},
		}

		result := mapper.ToAssetSummaryProtoList(domains)

		assert.Len(t, result, 2)
		assert.Equal(t, id1.String(), result[0].Id)
		assert.Equal(t, id2.String(), result[1].Id)
	})

	t.Run("Should return nil or empty when input is nil", func(t *testing.T) {
		result := mapper.ToAssetSummaryProtoList(nil)
		assert.Nil(t, result)
	})
}
