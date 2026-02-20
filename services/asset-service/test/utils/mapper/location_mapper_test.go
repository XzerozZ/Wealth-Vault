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

func TestToLocationDomain(t *testing.T) {
	tests := []struct {
		name     string
		input    *pb.Location
		expected domain.Location
	}{
		{
			name: "Success - Full Location Mapping",
			input: &pb.Location{
				Address:     "123 Sukhumvit Rd",
				Subdistrict: "Khlong Toei",
				District:    "Khlong Toei",
				Province:    "Bangkok",
				PostalCode:  "10110",
			},
			expected: domain.Location{
				Address:     "123 Sukhumvit Rd",
				Subdistrict: "Khlong Toei",
				District:    "Khlong Toei",
				Province:    "Bangkok",
				PostalCode:  "10110",
			},
		},
		{
			name:     "Edge Case - Nil Input should return empty object",
			input:    nil,
			expected: domain.Location{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToLocationDomain(tt.input)

			assert.Equal(t, tt.expected.Address, result.Address)
			assert.Equal(t, tt.expected.Subdistrict, result.Subdistrict)
			assert.Equal(t, tt.expected.District, result.District)
			assert.Equal(t, tt.expected.Province, result.Province)
			assert.Equal(t, tt.expected.PostalCode, result.PostalCode)
		})
	}
}

func TestToLocationProto(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	tests := []struct {
		name     string
		input    *domain.Location
		expected string
	}{
		{
			name: "Success - Map Domain to Proto",
			input: &domain.Location{
				ID:          id,
				Address:     "99/9 Ratchadaphisek",
				Subdistrict: "Din Daeng",
				CreatedAt:   now,
			},
			expected: "99/9 Ratchadaphisek",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToLocationProto(tt.input)

			assert.Equal(t, tt.input.ID.String(), result.Id)
			assert.Equal(t, tt.expected, result.Address)
			assert.Equal(t, tt.input.Subdistrict, result.Subdistrict)

			assert.NotNil(t, result.CreatedAt)
			assert.NotNil(t, result.UpdatedAt)
		})
	}
}
