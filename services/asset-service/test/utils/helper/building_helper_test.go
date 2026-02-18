package helper_test

import (
	"testing"

	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils/helper"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/protobuf/field_mask"
)

func TestApplyUpdateBuildingFields(t *testing.T) {
	tests := []struct {
		name     string
		req      *pb.UpdateBuildingRequest
		initial  *domain.Building
		expected *domain.Building
	}{
		{
			name: "Success - Update name and simple fields",
			req: &pb.UpdateBuildingRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"name", "amount", "type"}},
				Building: &pb.Building{
					Name:   "The Luxury Home",
					Amount: 12000000,
					Type:   pb.BuildingType_BUILDING_TYPE_HOUSE,
				},
			},
			initial: &domain.Building{
				Name:   "Old Home",
				Amount: 5000000,
				Type:   domain.BuildingTypeCondo,
			},
			expected: &domain.Building{
				Name:   "The Luxury Home",
				Amount: 12000000,
				Type:   domain.BuildingTypeHouse,
			},
		},
		{
			name: "Success - Update Nested Location fields",
			req: &pb.UpdateBuildingRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"location.address", "location.province"}},
				Building: &pb.Building{
					Location: &pb.Location{
						Address:  "99/1 Sukhumvit",
						Province: "Bangkok",
					},
				},
			},
			initial: &domain.Building{
				Location: domain.Location{
					Address:  "Original St.",
					Province: "Phuket",
					District: "Kathu",
				},
			},
			expected: &domain.Building{
				Location: domain.Location{
					Address:  "99/1 Sukhumvit",
					Province: "Bangkok",
					District: "Kathu",
				},
			},
		},
		{
			name: "Edge Case - Invalid Type should not update",
			req: &pb.UpdateBuildingRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{"type"}},
				Building: &pb.Building{
					Type: pb.BuildingType_BUILDING_TYPE_UNSPECIFIED,
				},
			},
			initial: &domain.Building{
				Type: domain.BuildingTypeTownHome,
			},
			expected: &domain.Building{
				Type: domain.BuildingTypeTownHome,
			},
		},
		{
			name: "Edge Case - Missing paths in UpdateMask",
			req: &pb.UpdateBuildingRequest{
				UpdateMask: &field_mask.FieldMask{Paths: []string{}},
				Building: &pb.Building{
					Name: "Changed but no mask",
				},
			},
			initial: &domain.Building{
				Name: "Original",
			},
			expected: &domain.Building{
				Name: "Original",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := helper.ApplyUpdateBuildingFields(tt.req, tt.initial)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Name, tt.initial.Name)
			assert.Equal(t, tt.expected.Amount, tt.initial.Amount)
			assert.Equal(t, tt.expected.Type, tt.initial.Type)
			assert.Equal(t, tt.expected.Location.Address, tt.initial.Location.Address)
			assert.Equal(t, tt.expected.Location.Province, tt.initial.Location.Province)
			assert.Equal(t, tt.expected.Location.District, tt.initial.Location.District)
		})
	}
}
