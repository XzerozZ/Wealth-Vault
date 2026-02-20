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

func TestToInsuranceDomain(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name     string
		input    *pb.CreateInsuranceRequest
		expected *domain.Insurance
	}{
		{
			name: "Success - Full Mapping",
			input: &pb.CreateInsuranceRequest{
				Name:           "Life Insurance A",
				Type:           pb.InsuranceType_INSURANCE_TYPE_LIFE,
				PolNum:         "POL123",
				CompanyName:    "Allianz",
				CoverageAmount: 1000000,
			},
			expected: &domain.Insurance{
				UserID:         userID,
				Name:           "Life Insurance A",
				Type:           domain.InsuranceTypeLife,
				PolicyNumber:   "POL123",
				CompanyName:    "Allianz",
				CoverageAmount: 1000000,
			},
		},
		{
			name: "Fallback - Unknown Type should use default",
			input: &pb.CreateInsuranceRequest{
				Type: pb.InsuranceType(999),
			},
			expected: &domain.Insurance{
				UserID: userID,
				Type:   domain.InsuranceTypeLife,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToInsuranceDomain(tt.input, userID)

			assert.Equal(t, tt.expected.UserID, result.UserID)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Type, result.Type)
			assert.Equal(t, tt.expected.PolicyNumber, result.PolicyNumber)
			assert.Equal(t, tt.expected.CoverageAmount, result.CoverageAmount)
		})
	}
}

func TestToInsuranceProto(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	tests := []struct {
		name     string
		input    *domain.Insurance
		expected pb.InsuranceType
	}{
		{
			name: "Map HEALTH correctly with dates",
			input: &domain.Insurance{
				ID:        id,
				UserID:    userID,
				Type:      "HEALTH",
				Name:      "Bupa Health",
				ConDate:   &now,
				CreatedAt: now,
			},
			expected: pb.InsuranceType_INSURANCE_TYPE_HEALTH,
		},
		{
			name: "Map Unknown Type and Nil Dates",
			input: &domain.Insurance{
				Type:    "UNKNOWN_INS",
				ConDate: nil,
				ExpDate: nil,
			},
			expected: pb.InsuranceType_INSURANCE_TYPE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToInsuranceProto(tt.input)

			assert.Equal(t, tt.expected, result.Type)

			if tt.input.ConDate != nil {
				assert.NotNil(t, result.ConDate)
			} else {
				assert.Nil(t, result.ConDate)
			}
		})
	}
}

func TestToInsuranceProtoSlice(t *testing.T) {
	t.Run("Empty slice should return empty slice not nil", func(t *testing.T) {
		result := mapper.ToInsuranceProtoSlice([]*domain.Insurance{})
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("Should map all elements", func(t *testing.T) {
		insurances := []*domain.Insurance{
			{Name: "Ins 1", CreatedAt: time.Now()},
			{Name: "Ins 2", CreatedAt: time.Now()},
		}
		result := mapper.ToInsuranceProtoSlice(insurances)
		assert.Len(t, result, 2)
		assert.Equal(t, "Ins 1", result[0].Name)
	})
}
