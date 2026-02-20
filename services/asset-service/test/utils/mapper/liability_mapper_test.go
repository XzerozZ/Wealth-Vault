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

func TestToLiabilityDomain(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name     string
		input    *pb.CreateLiabilityRequest
		expected *domain.Liability
	}{
		{
			name: "Success - Valid Loan Mapping",
			input: &pb.CreateLiabilityRequest{
				Name:         "Home Loan",
				Type:         pb.LiabilityType_LIABILITY_TYPE_LOAN,
				Creditor:     "SCB Bank",
				Principal:    5000000,
				InterestRate: 2.5,
				Description:  "Housing finance",
			},
			expected: &domain.Liability{
				UserID:       userID,
				Name:         "Home Loan",
				Type:         domain.LiabilityTypeLoan,
				Creditor:     "SCB Bank",
				Principal:    5000000,
				InterestRate: 2.5,
				Description:  "Housing finance",
			},
		},
		{
			name: "Fallback - Unknown Liability Type",
			input: &pb.CreateLiabilityRequest{
				Type: pb.LiabilityType(999),
			},
			expected: &domain.Liability{
				UserID: userID,
				Type:   domain.LiabilityTypeLoan,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToLiabilityDomain(tt.input, userID)

			assert.Equal(t, tt.expected.UserID, result.UserID)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Type, result.Type)
			assert.Equal(t, tt.expected.Principal, result.Principal)
			assert.Equal(t, tt.expected.InterestRate, result.InterestRate)
			assert.Equal(t, tt.expected.Creditor, result.Creditor)
		})
	}
}

func TestToLiabilityProto(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	tests := []struct {
		name     string
		input    *domain.Liability
		expected pb.LiabilityType
	}{
		{
			name: "Map EXPENSE correctly",
			input: &domain.Liability{
				ID:        id,
				UserID:    userID,
				Type:      "EXPENSE",
				Name:      "Credit Card",
				StartAt:   &now,
				CreatedAt: now,
			},
			expected: pb.LiabilityType_LIABILITY_TYPE_EXPENSE,
		},
		{
			name: "Map Loan Case-Insensitive with spaces",
			input: &domain.Liability{
				Type: " loan ",
			},
			expected: pb.LiabilityType_LIABILITY_TYPE_LOAN,
		},
		{
			name: "Map Unknown Type to UNSPECIFIED",
			input: &domain.Liability{
				Type: "GAMBLING_DEBT",
			},
			expected: pb.LiabilityType_LIABILITY_TYPE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.CreatedAt = now
			result := mapper.ToLiabilityProto(tt.input)

			assert.Equal(t, tt.expected, result.Type)
			if tt.input.ID != uuid.Nil {
				assert.Equal(t, tt.input.ID.String(), result.Id)
			}

			if tt.input.StartAt != nil {
				assert.NotNil(t, result.StartAt)
			}
		})
	}
}

func TestToLiabilityProtoSlice(t *testing.T) {
	t.Run("Empty slice should return empty slice not nil", func(t *testing.T) {
		result := mapper.ToLiabilityProtoSlice([]*domain.Liability{})
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("Should map all liability elements", func(t *testing.T) {
		liabilities := []*domain.Liability{
			{Name: "Debt 1", CreatedAt: time.Now()},
			{Name: "Debt 2", CreatedAt: time.Now()},
		}
		result := mapper.ToLiabilityProtoSlice(liabilities)
		assert.Len(t, result, 2)
		assert.Equal(t, "Debt 1", result[0].Name)
	})
}
