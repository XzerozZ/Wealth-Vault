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

func TestToCashDomain(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name     string
		input    *pb.CreateCashRequest
		expected *domain.Cash
	}{
		{
			name: "Success - Valid Mapping",
			input: &pb.CreateCashRequest{
				Name:        "Emergency Cash",
				Amount:      5000.0,
				Description: "Cash in drawer",
			},
			expected: &domain.Cash{
				UserID:      userID,
				Name:        "Emergency Cash",
				Amount:      5000.0,
				Description: "Cash in drawer",
			},
		},
		{
			name: "Success - Minimal Data",
			input: &pb.CreateCashRequest{
				Name:   "Pocket Money",
				Amount: 100.0,
			},
			expected: &domain.Cash{
				UserID: userID,
				Name:   "Pocket Money",
				Amount: 100.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToCashDomain(tt.input, userID)

			assert.Equal(t, tt.expected.UserID, result.UserID)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Amount, result.Amount)
			assert.Equal(t, tt.expected.Description, result.Description)
		})
	}
}

func TestToCashProto(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	tests := []struct {
		name  string
		input *domain.Cash
	}{
		{
			name: "Success - Full Mapping",
			input: &domain.Cash{
				ID:          id,
				UserID:      userID,
				Name:        "Savings Cash",
				Amount:      2500.75,
				Description: "In safe box",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToCashProto(tt.input)

			assert.Equal(t, tt.input.ID.String(), result.Id)
			assert.Equal(t, tt.input.UserID.String(), result.UserId)
			assert.Equal(t, tt.input.Name, result.Name)
			assert.Equal(t, tt.input.Amount, result.Amount)
			assert.Equal(t, tt.input.Description, result.Description)

			assert.NotNil(t, result.CreatedAt)
			assert.NotNil(t, result.UpdatedAt)
		})
	}
}

func TestToCashProtoSlice(t *testing.T) {
	t.Run("Empty slice should return empty slice not nil", func(t *testing.T) {
		result := mapper.ToCashProtoSlice([]*domain.Cash{})
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("Should map all elements correctly", func(t *testing.T) {
		cashes := []*domain.Cash{
			{ID: uuid.New(), Name: "Cash 1", CreatedAt: time.Now()},
			{ID: uuid.New(), Name: "Cash 2", CreatedAt: time.Now()},
		}
		result := mapper.ToCashProtoSlice(cashes)

		assert.Len(t, result, 2)
		assert.Equal(t, cashes[0].Name, result[0].Name)
		assert.Equal(t, cashes[1].Name, result[1].Name)
	})
}
