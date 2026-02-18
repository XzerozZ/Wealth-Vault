package helper_test

import (
	"testing"
	"time"

	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils/helper"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestApplyUpdateLiabilityFields(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	protoTime := timestamppb.New(now)

	tests := []struct {
		name     string
		req      *pb.UpdateLiabilityRequest
		initial  *domain.Liability
		expected *domain.Liability
	}{
		{
			name: "Success - Update financial fields and dates",
			req: &pb.UpdateLiabilityRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"principal", "interest_rate", "started_at"}},
				Liability: &pb.Liability{
					Principal:    1000000,
					InterestRate: 5.5,
					StartAt:      protoTime,
				},
			},
			initial: &domain.Liability{
				Principal:    500000,
				InterestRate: 3.0,
			},
			expected: &domain.Liability{
				Principal:    1000000,
				InterestRate: 5.5,
				StartAt:      &now,
			},
		},
		{
			name: "Success - Update type to EXPENSE",
			req: &pb.UpdateLiabilityRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"type"}},
				Liability: &pb.Liability{
					Type: pb.LiabilityType_LIABILITY_TYPE_EXPENSE,
				},
			},
			initial: &domain.Liability{
				Type: domain.LiabilityTypeLoan,
			},
			expected: &domain.Liability{
				Type: domain.LiabilityTypeExpense,
			},
		},
		{
			name: "Check Bug - Creditor should update correctly",
			req: &pb.UpdateLiabilityRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"creditor"}},
				Liability: &pb.Liability{
					Creditor: "K-Bank",
				},
			},
			initial: &domain.Liability{
				Name:     "Personal Loan",
				Creditor: "Old Bank",
			},
			expected: &domain.Liability{
				Name:     "Personal Loan",
				Creditor: "K-Bank",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := helper.ApplyUpdateFields(tt.req, tt.initial)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Principal, tt.initial.Principal)
			assert.Equal(t, tt.expected.InterestRate, tt.initial.InterestRate)
			assert.Equal(t, tt.expected.Creditor, tt.initial.Creditor)
			assert.Equal(t, tt.expected.Name, tt.initial.Name)
			assert.Equal(t, tt.expected.Type, tt.initial.Type)

			if tt.expected.StartAt != nil {
				assert.Equal(t, tt.expected.StartAt.Unix(), tt.initial.StartAt.Unix())
			}
		})
	}
}
