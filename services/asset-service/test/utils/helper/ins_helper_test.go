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

func TestApplyUpdateInsuranceFields(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	protoTime := timestamppb.New(now)

	tests := []struct {
		name     string
		req      *pb.UpdateInsuranceRequest
		initial  *domain.Insurance
		expected *domain.Insurance
	}{
		{
			name: "Success - Update policy details and dates",
			req: &pb.UpdateInsuranceRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"policy_number", "con_date", "exp_date"}},
				Insurance: &pb.Insurance{
					PolNum:  "NEW-POL-999",
					ConDate: protoTime,
					ExpDate: protoTime,
				},
			},
			initial: &domain.Insurance{
				PolicyNumber: "OLD-123",
				ConDate:      nil,
			},
			expected: &domain.Insurance{
				PolicyNumber: "NEW-POL-999",
				ConDate:      &now,
				ExpDate:      &now,
			},
		},
		{
			name: "Success - Update insurance type",
			req: &pb.UpdateInsuranceRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"type"}},
				Insurance: &pb.Insurance{
					Type: pb.InsuranceType_INSURANCE_TYPE_HEALTH,
				},
			},
			initial: &domain.Insurance{
				Type: domain.InsuranceTypeLife,
			},
			expected: &domain.Insurance{
				Type: domain.InsuranceTypeHealth,
			},
		},
		{
			name: "Edge Case - Null Date in Request (Should skip)",
			req: &pb.UpdateInsuranceRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"con_date"}},
				Insurance: &pb.Insurance{
					ConDate: nil,
				},
			},
			initial: &domain.Insurance{
				Name:    "Life Plan",
				ConDate: &now,
			},
			expected: &domain.Insurance{
				Name:    "Life Plan",
				ConDate: &now,
			},
		},
		{
			name: "Edge Case - Unspecified Type (Should skip)",
			req: &pb.UpdateInsuranceRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"type"}},
				Insurance: &pb.Insurance{
					Type: pb.InsuranceType_INSURANCE_TYPE_UNSPECIFIED,
				},
			},
			initial: &domain.Insurance{
				Type: domain.InsuranceTypeAccident,
			},
			expected: &domain.Insurance{
				Type: domain.InsuranceTypeAccident,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := helper.ApplyUpdateInsuranceFields(tt.req, tt.initial)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.PolicyNumber, tt.initial.PolicyNumber)
			assert.Equal(t, tt.expected.Type, tt.initial.Type)

			if tt.expected.ConDate != nil {
				assert.NotNil(t, tt.initial.ConDate)
				assert.Equal(t, tt.expected.ConDate.Unix(), tt.initial.ConDate.Unix())
			}
			if tt.expected.ExpDate != nil {
				assert.NotNil(t, tt.initial.ExpDate)
				assert.Equal(t, tt.expected.ExpDate.Unix(), tt.initial.ExpDate.Unix())
			}
		})
	}
}
