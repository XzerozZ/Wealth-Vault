package mapper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"
	"wealth-vault/user-service/pkg/utils/mapper"
)

func TestAssetMapper(t *testing.T) {
	now := time.Now()
	recentDeletedAt := timestamppb.New(now.Add(-24 * time.Hour))
	oldDeletedAt := timestamppb.New(now.Add(-10 * 24 * time.Hour))

	t.Run("MapBuilding - Active Case with Location", func(t *testing.T) {
		building := &assetPb.Building{
			Id:     "b1",
			Name:   "Condo A",
			Amount: 5000000,
			Location: &assetPb.Location{
				District: "ปทุมวัน",
				Province: "กรุงเทพ",
			},
			Type: assetPb.BuildingType_BUILDING_TYPE_CONDO,
		}

		result := mapper.MapBuildingToPreview(building)

		assert.NotNil(t, result)
		preview := result.GetBuilding()
		assert.Equal(t, "b1", preview.Id)
		assert.Equal(t, "ปทุมวัน, กรุงเทพ", preview.LocationText)
		assert.Equal(t, "BUILDING_TYPE_CONDO", preview.TypeName)
	})

	t.Run("MapBuilding - Ghost Case (Recently Deleted)", func(t *testing.T) {
		building := &assetPb.Building{
			Id:        "b1",
			Name:      "Condo A",
			DeletedAt: recentDeletedAt,
		}

		result := mapper.MapBuildingToPreview(building)

		assert.NotNil(t, result)
		deleted := result.GetDeleted()
		assert.NotNil(t, deleted)
		assert.Equal(t, "Condo A", deleted.OriginalName)
		assert.Contains(t, deleted.Message, "ถูกลบไปแล้ว")
	})

	t.Run("MapBuilding - Expired Case (Deleted too long)", func(t *testing.T) {
		building := &assetPb.Building{
			Id:        "b1",
			DeletedAt: oldDeletedAt,
		}

		result := mapper.MapBuildingToPreview(building)

		assert.Nil(t, result, "Should be nil if deleted more than 7 days")
	})

	t.Run("MapBuilding - Nil Location Safety", func(t *testing.T) {
		building := &assetPb.Building{
			Id:       "b1",
			Location: nil,
		}

		result := mapper.MapBuildingToPreview(building)
		assert.Equal(t, "ไม่ระบุตำแหน่ง", result.GetBuilding().LocationText)
	})

	t.Run("MapAccount - Normal Case with Masking", func(t *testing.T) {
		account := &assetPb.Account{
			Id:      "acc1",
			Name:    "Savings",
			BankAcc: "1234567890",
		}

		result := mapper.MapAccountToPreview(account)

		assert.NotNil(t, result)
		accPreview := result.GetAccount()
		assert.Contains(t, accPreview.AccountNumber, "7890")
	})

	t.Run("MapInsurance - Nil ExpDate Safety", func(t *testing.T) {
		insurance := &assetPb.Insurance{
			Id:      "ins1",
			ExpDate: nil,
		}

		result := mapper.MapInsuranceToPreview(insurance)
		assert.Equal(t, "ไม่ระบุ", result.GetInsurance().ExpDateText)
	})

	t.Run("Input Nil - Should return Nil", func(t *testing.T) {
		assert.Nil(t, mapper.MapBuildingToPreview(nil))
		assert.Nil(t, mapper.MapAccountToPreview(nil))
	})
}
