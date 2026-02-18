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

func TestMapToLands(t *testing.T) {
	ids := []uuid.UUID{uuid.New(), uuid.New()}

	t.Run("Should map UUID slice to Land domain slice correctly", func(t *testing.T) {
		result := mapper.MapToLands(ids)

		assert.Len(t, result, 2)
		assert.Equal(t, ids[0], result[0].ID)
		assert.Equal(t, ids[1], result[1].ID)
	})

	t.Run("Empty slice should return empty result", func(t *testing.T) {
		result := mapper.MapToLands([]uuid.UUID{})
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})
}

func TestToLandDomain(t *testing.T) {
	userID := uuid.New()
	buildingID := uuid.New()

	tests := []struct {
		name     string
		input    *pb.CreateLandRequest
		expected *domain.Land
	}{
		{
			name: "Success - Valid Mapping",
			input: &pb.CreateLandRequest{
				Name:        "Green Field",
				DeedNum:     "SN-12345",
				Area:        50,
				Amount:      1500000,
				BuildingIds: []string{buildingID.String()},
			},
			expected: &domain.Land{
				UserID:  userID,
				Name:    "Green Field",
				DeedNum: "SN-12345",
				Area:    50,
				Amount:  1500000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToLandDomain(tt.input, userID)

			assert.Equal(t, tt.expected.UserID, result.UserID)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.DeedNum, result.DeedNum)
			assert.Equal(t, tt.expected.Amount, result.Amount)

			// เช็คความสัมพันธ์ว่า Building ID ถูกแปลงมาใส่ใน Lands จริงไหม
			if len(tt.input.BuildingIds) > 0 {
				assert.NotEmpty(t, result.Buildings)
				assert.Equal(t, buildingID, result.Buildings[0].ID)
			}
		})
	}
}

func TestToLandProto(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	t.Run("Success - Full Mapping to Proto", func(t *testing.T) {
		d := &domain.Land{
			ID:        id,
			UserID:    userID,
			Name:      "Beach Land",
			DeedNum:   "B-001",
			CreatedAt: now,
		}

		result := mapper.ToLandProto(d)

		assert.Equal(t, id.String(), result.Id)
		assert.Equal(t, userID.String(), result.UserId)
		assert.Equal(t, "Beach Land", result.Name)
		assert.Equal(t, "B-001", result.DeedNum)
		assert.NotNil(t, result.CreatedAt)
	})
}

func TestToLandProtoSlice(t *testing.T) {
	t.Run("Empty slice should return empty slice not nil", func(t *testing.T) {
		result := mapper.ToLandProtoSlice([]*domain.Land{})
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("Should map all land elements correctly", func(t *testing.T) {
		lands := []*domain.Land{
			{Name: "Land A", CreatedAt: time.Now()},
			{Name: "Land B", CreatedAt: time.Now()},
		}
		result := mapper.ToLandProtoSlice(lands)
		assert.Len(t, result, 2)
		assert.Equal(t, "Land A", result[0].Name)
	})
}
