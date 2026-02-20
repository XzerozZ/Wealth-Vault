package helper_test

import (
	"testing"

	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils/helper"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/protobuf/field_mask"
)

func TestApplyUpdateInFields(t *testing.T) {
	tests := []struct {
		name     string
		req      *pb.UpdateInvestmentRequest
		initial  *domain.Investment
		expected *domain.Investment
	}{
		{
			name: "Success - Update symbol and quantity",
			req: &pb.UpdateInvestmentRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"symbol", "quantity"}},
				Invest: &pb.Investment{
					Symbol:   "BTC",
					Quantity: 0.5,
				},
			},
			initial: &domain.Investment{
				Symbol:   "ETH",
				Quantity: 2.0,
			},
			expected: &domain.Investment{
				Symbol:   "BTC",
				Quantity: 0.5,
			},
		},
		{
			name: "Success - Update investment type",
			req: &pb.UpdateInvestmentRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"type"}},
				Invest: &pb.Investment{
					Type: pb.InvestmentType_INVEST_TYPE_GOLD,
				},
			},
			initial: &domain.Investment{
				Type: domain.InvestTypeStockTH,
			},
			expected: &domain.Investment{
				Type: domain.InvestTypeGold,
			},
		},
		{
			name: "Edge Case - Numeric field is 0 in request (Should skip)",
			req: &pb.UpdateInvestmentRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"cost_per_price"}},
				Invest: &pb.Investment{
					CostPrice: 0,
				},
			},
			initial: &domain.Investment{
				CostPerPrice: 150.25,
			},
			expected: &domain.Investment{
				CostPerPrice: 150.25,
			},
		},
		{
			name: "Edge Case - Unspecified Type (Should skip)",
			req: &pb.UpdateInvestmentRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"type"}},
				Invest: &pb.Investment{
					Type: pb.InvestmentType_INVEST_TYPE_UNSPECIFIED,
				},
			},
			initial: &domain.Investment{
				Type: domain.InvestTypeCrypto,
			},
			expected: &domain.Investment{
				Type: domain.InvestTypeCrypto,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := helper.ApplyUpdateInFields(tt.req, tt.initial)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Symbol, tt.initial.Symbol)
			assert.Equal(t, tt.expected.Quantity, tt.initial.Quantity)
			assert.Equal(t, tt.expected.Type, tt.initial.Type)
			assert.Equal(t, tt.expected.CostPerPrice, tt.initial.CostPerPrice)
		})
	}
}
