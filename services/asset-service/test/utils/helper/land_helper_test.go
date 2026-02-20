package helper_test

import (
	"testing"

	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils/helper"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/protobuf/field_mask"
)

func TestApplyUpdateLandFields(t *testing.T) {
	tests := []struct {
		name     string
		req      *pb.UpdateLandRequest
		initial  *domain.Land
		expected *domain.Land
	}{
		{
			name: "Success - Update basic land fields",
			req: &pb.UpdateLandRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"name", "deed_num", "amount"}},
				Land: &pb.Land{
					Name:    "New Orchard",
					DeedNum: "NS3K-1234",
					Amount:  2500000,
				},
			},
			initial: &domain.Land{
				Name:    "Old Plot",
				DeedNum: "9999",
				Amount:  1000000,
			},
			expected: &domain.Land{
				Name:    "New Orchard",
				DeedNum: "NS3K-1234",
				Amount:  2500000,
			},
		},
		{
			name: "Success - Update Nested Location",
			req: &pb.UpdateLandRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"location.address", "location.district"}},
				Land: &pb.Land{
					Location: &pb.Location{
						Address:  "123 New Road",
						District: "Bangkok Noi",
					},
				},
			},
			initial: &domain.Land{
				Location: domain.Location{
					Address:  "Old Alley",
					District: "Bangkok Yai",
					Province: "Bangkok",
				},
			},
			expected: &domain.Land{
				Location: domain.Location{
					Address:  "123 New Road",
					District: "Bangkok Noi",
					Province: "Bangkok",
				},
			},
		},
		{
			name: "Edge Case - Numeric area is 0 in request (Should skip)",
			req: &pb.UpdateLandRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"area"}},
				Land: &pb.Land{
					Area: 0,
				},
			},
			initial: &domain.Land{
				Area: 150.5,
			},
			expected: &domain.Land{
				Area: 150.5,
			},
		},
		{
			name: "Edge Case - Path not in mask (Should not update)",
			req: &pb.UpdateLandRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"name"}},
				Land: &pb.Land{
					Name:        "New Name",
					Description: "Should not change",
				},
			},
			initial: &domain.Land{
				Name:        "Old Name",
				Description: "Original desc",
			},
			expected: &domain.Land{
				Name:        "New Name",
				Description: "Original desc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := helper.ApplyUpdateLandFields(tt.req, tt.initial)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Name, tt.initial.Name)
			assert.Equal(t, tt.expected.DeedNum, tt.initial.DeedNum)
			assert.Equal(t, tt.expected.Amount, tt.initial.Amount)
			assert.Equal(t, tt.expected.Location.Address, tt.initial.Location.Address)
			assert.Equal(t, tt.expected.Location.District, tt.initial.Location.District)
			assert.Equal(t, tt.expected.Location.Province, tt.initial.Location.Province)
		})
	}
}
