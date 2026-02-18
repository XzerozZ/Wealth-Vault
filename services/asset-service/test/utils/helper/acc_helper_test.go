package helper_test

import (
	"testing"

	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils/helper"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/protobuf/field_mask"
)

func TestApplyUpdateAccFields(t *testing.T) {
	tests := []struct {
		name     string
		req      *pb.UpdateAccountRequest
		initial  *domain.Account
		expected *domain.Account
	}{
		{
			name: "Success - Update name and amount",
			req: &pb.UpdateAccountRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"name", "amount"}},
				Acc: &pb.Account{
					Name:   "New Name",
					Amount: 5000.0,
				},
			},
			initial: &domain.Account{
				Name:   "Old Name",
				Amount: 1000.0,
			},
			expected: &domain.Account{
				Name:   "New Name",
				Amount: 5000.0,
			},
		},
		{
			name: "Success - Update type",
			req: &pb.UpdateAccountRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"type"}},
				Acc: &pb.Account{
					Type: pb.BankAccType_BANK_ACC_TYPE_FIXED_DEPOSIT,
				},
			},
			initial: &domain.Account{
				Type: domain.BankTypeSavings,
			},
			expected: &domain.Account{
				Type: domain.BankTypeFixedDeposit,
			},
		},
		{
			name: "Partial Skip - Path exists but value is empty (should not update)",
			req: &pb.UpdateAccountRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"description"}},
				Acc: &pb.Account{
					Description: "",
				},
			},
			initial: &domain.Account{
				Description: "Keep this",
			},
			expected: &domain.Account{
				Description: "Keep this",
			},
		},
		{
			name: "Empty Mask - No fields should change",
			req: &pb.UpdateAccountRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{}},
				Acc: &pb.Account{
					Name: "New Name",
				},
			},
			initial: &domain.Account{
				Name: "Old Name",
			},
			expected: &domain.Account{
				Name: "Old Name",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := helper.ApplyUpdateAccFields(tt.req, tt.initial)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Name, tt.initial.Name)
			assert.Equal(t, tt.expected.Amount, tt.initial.Amount)
			assert.Equal(t, tt.expected.Type, tt.initial.Type)
			assert.Equal(t, tt.expected.Description, tt.initial.Description)
		})
	}
}
