package helper_test

import (
	"testing"

	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils/helper"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/protobuf/field_mask"
)

func TestApplyUpdateCashFields(t *testing.T) {
	tests := []struct {
		name     string
		req      *pb.UpdateCashRequest
		initial  *domain.Cash
		expected *domain.Cash
	}{
		{
			name: "Success - Update name and amount",
			req: &pb.UpdateCashRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"name", "amount"}},
				Cash: &pb.Cash{
					Name:   "Emergency Fund",
					Amount: 15000.75,
				},
			},
			initial: &domain.Cash{
				Name:   "Old Cash",
				Amount: 5000.0,
			},
			expected: &domain.Cash{
				Name:   "Emergency Fund",
				Amount: 15000.75,
			},
		},
		{
			name: "Success - Update only description",
			req: &pb.UpdateCashRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"description"}},
				Cash: &pb.Cash{
					Description: "Updated description",
				},
			},
			initial: &domain.Cash{
				Name:        "My Wallet",
				Description: "Original desc",
			},
			expected: &domain.Cash{
				Name:        "My Wallet",
				Description: "Updated description",
			},
		},
		{
			name: "Edge Case - Path in mask but zero value in request",
			req: &pb.UpdateCashRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"amount"}},
				Cash: &pb.Cash{
					Amount: 0, // ค่า 0 ไม่ควรทับค่าเดิมตาม logic ใน switch-case
				},
			},
			initial: &domain.Cash{
				Amount: 99.0,
			},
			expected: &domain.Cash{
				Amount: 99.0,
			},
		},
		{
			name: "Edge Case - Empty UpdateMask",
			req: &pb.UpdateCashRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{}},
				Cash: &pb.Cash{
					Name: "Try to change",
				},
			},
			initial: &domain.Cash{
				Name: "Original",
			},
			expected: &domain.Cash{
				Name: "Original",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := helper.ApplyUpdateCashFields(tt.req, tt.initial)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Name, tt.initial.Name)
			assert.Equal(t, tt.expected.Amount, tt.initial.Amount)
			assert.Equal(t, tt.expected.Description, tt.initial.Description)
		})
	}
}
