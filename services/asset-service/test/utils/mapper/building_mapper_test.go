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

func TestMapToBuildings(t *testing.T) {
	ids := []uuid.UUID{uuid.New(), uuid.New()}

	t.Run("Should map UUID slice to Building domain slice correctly", func(t *testing.T) {
		result := mapper.MapToBuildings(ids)

		assert.Len(t, result, 2)
		assert.Equal(t, ids[0], result[0].ID)
		assert.Equal(t, ids[1], result[1].ID)
	})

	t.Run("Empty slice should return empty result", func(t *testing.T) {
		result := mapper.MapToBuildings([]uuid.UUID{})
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})
}

func TestToBuildingDomain(t *testing.T) {
	userID := uuid.New()
	landID := uuid.New()

	tests := []struct {
		name     string
		input    *pb.CreateBuildingRequest
		expected *domain.Building
	}{
		{
			name: "Success - Valid Mapping",
			input: &pb.CreateBuildingRequest{
				Name:        "The Grand Condo",
				Type:        pb.BuildingType_BUILDING_TYPE_CONDO,
				Area:        50.5,
				Amount:      5000000,
				Description: "Luxury condo",
				LandIds:     []string{landID.String()},
			},
			expected: &domain.Building{
				UserID:      userID,
				Name:        "The Grand Condo",
				Type:        domain.BuildingTypeCondo,
				Area:        50.5,
				Amount:      5000000,
				Description: "Luxury condo",
			},
		},
		{
			name: "Fallback - Unknown Type should use default",
			input: &pb.CreateBuildingRequest{
				Type: pb.BuildingType(999),
			},
			expected: &domain.Building{
				UserID: userID,
				Type:   domain.BuildingTypeHouse,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToBuildingDomain(tt.input, userID)

			assert.Equal(t, tt.expected.UserID, result.UserID)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Type, result.Type)
			assert.Equal(t, tt.expected.Area, result.Area)
			assert.Equal(t, tt.expected.Amount, result.Amount)

			if len(tt.input.LandIds) > 0 {
				assert.NotEmpty(t, result.Lands)
				assert.Equal(t, tt.input.LandIds[0], result.Lands[0].ID.String())
			}
		})
	}
}

func TestToBuildingProto(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	tests := []struct {
		name     string
		input    *domain.Building
		expected pb.BuildingType
	}{
		{
			name: "Map CONDO correctly",
			input: &domain.Building{
				ID:     id,
				UserID: userID,
				Type:   "CONDO",
			},
			expected: pb.BuildingType_BUILDING_TYPE_CONDO,
		},
		{
			name: "Map Case-Insensitive with spaces",
			input: &domain.Building{
				Type: " townhome ",
			},
			expected: pb.BuildingType_BUILDING_TYPE_TOWNHOME,
		},
		{
			name: "Map Unknown Type to UNSPECIFIED",
			input: &domain.Building{
				Type: "PENTHOUSE",
			},
			expected: pb.BuildingType_BUILDING_TYPE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.CreatedAt = now
			result := mapper.ToBuildingProto(tt.input)

			assert.Equal(t, tt.expected, result.Type)
			if tt.input.ID != uuid.Nil {
				assert.Equal(t, tt.input.ID.String(), result.Id)
			}
		})
	}
}

func TestToBuildingProtoSlice(t *testing.T) {
	t.Run("Empty slice should return empty slice not nil", func(t *testing.T) {
		result := mapper.ToBuildingProtoSlice([]*domain.Building{})
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("Should map all elements", func(t *testing.T) {
		buildings := []*domain.Building{
			{Name: "House A", CreatedAt: time.Now()},
			{Name: "House B", CreatedAt: time.Now()},
		}
		result := mapper.ToBuildingProtoSlice(buildings)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, "House A", result[0].Name)
		assert.Equal(t, "House B", result[1].Name)
	})
}
