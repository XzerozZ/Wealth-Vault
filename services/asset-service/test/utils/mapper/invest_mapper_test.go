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

func TestToInvestmentDomain(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name     string
		input    *pb.CreateInvestmentRequest
		expected *domain.Investment
	}{
		{
			name: "Success - Valid Mapping (Stock)",
			input: &pb.CreateInvestmentRequest{
				Name:        "Apple Inc.",
				Symbol:      "AAPL",
				Type:        pb.InvestmentType_INVEST_TYPE_STOCK_US,
				BrokerName:  "Interactive Brokers",
				Quantity:    10.5,
				CostPrice:   150.0,
				Description: "Long term investment",
			},
			expected: &domain.Investment{
				UserID:       userID,
				Name:         "Apple Inc.",
				Symbol:       "AAPL",
				Type:         domain.InvestTypeStockUS,
				BrokerName:   "Interactive Brokers",
				Quantity:     10.5,
				CostPerPrice: 150.0,
				Description:  "Long term investment",
			},
		},
		{
			name: "Fallback - Unknown Investment Type",
			input: &pb.CreateInvestmentRequest{
				Type: pb.InvestmentType(999),
			},
			expected: &domain.Investment{
				UserID: userID,
				Type:   domain.InvestTypeStockUS,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.ToInvestmentDomain(tt.input, userID)

			assert.Equal(t, tt.expected.UserID, result.UserID)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Type, result.Type)
			assert.Equal(t, tt.expected.Quantity, result.Quantity)
			assert.Equal(t, tt.expected.CostPerPrice, result.CostPerPrice)
		})
	}
}

func TestToInvestProto(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	tests := []struct {
		name     string
		input    *domain.Investment
		expected pb.InvestmentType
	}{
		{
			name: "Map CRYPTO correctly",
			input: &domain.Investment{
				ID:        id,
				UserID:    userID,
				Type:      "CRYPTO",
				Symbol:    "BTC",
				CreatedAt: now,
			},
			expected: pb.InvestmentType_INVEST_TYPE_CRYPTO,
		},
		{
			name: "Map Mutual Fund Case-Insensitive",
			input: &domain.Investment{
				Type: " mutual_fund ",
			},
			expected: pb.InvestmentType_INVEST_TYPE_MUTUAL_FUND,
		},
		{
			name: "Map Unknown Type to UNSPECIFIED",
			input: &domain.Investment{
				Type: "REAL_ESTATE_TOKEN",
			},
			expected: pb.InvestmentType_INVEST_TYPE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.CreatedAt = now
			result := mapper.ToInvestProto(tt.input)

			assert.Equal(t, tt.expected, result.Type)
			if tt.input.Symbol != "" {
				assert.Equal(t, tt.input.Symbol, result.Symbol)
			}
		})
	}
}

func TestToInvestProtoSlice(t *testing.T) {
	t.Run("Empty slice should return empty slice not nil", func(t *testing.T) {
		result := mapper.ToInvestProtoSlice([]*domain.Investment{})
		assert.NotNil(t, result)
		assert.Empty(t, result)
	})

	t.Run("Should map all investment elements", func(t *testing.T) {
		invests := []*domain.Investment{
			{Name: "Stock A", CreatedAt: time.Now()},
			{Name: "Stock B", CreatedAt: time.Now()},
		}
		result := mapper.ToInvestProtoSlice(invests)
		assert.Len(t, result, 2)
		assert.Equal(t, "Stock A", result[0].Name)
	})
}
