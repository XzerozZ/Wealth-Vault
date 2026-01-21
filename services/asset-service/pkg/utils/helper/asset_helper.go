package helper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"

	"gorm.io/datatypes"
)

func ApplyUpdateAssetFields(req *pb.UpdateAssetRequest, asset *domain.Asset) ([]string, error) {
	var updateMask []string
	has := func(target string) bool {
		if req.UpdateMask == nil || len(req.UpdateMask.Paths) == 0 {
			return true
		}

		for _, p := range req.UpdateMask.Paths {
			if p == target {
				return true
			}
		}

		return false
	}

	if has("name") {
		asset.Name = req.Name
		updateMask = append(updateMask, "Name")
	}
	if has("amount") {
		asset.Amount = req.Amount
		updateMask = append(updateMask, "Amount")
	}
	if has("description") {
		asset.Description = req.Description
		updateMask = append(updateMask, "Description")
	}
	if has("detail") && req.Detail != nil {
		jsonBytes, err := MapProtoDetailToJSON(req.Detail)
		if err != nil {
			return nil, err
		}

		asset.Details = datatypes.JSON(jsonBytes)
		updateMask = append(updateMask, "Details")
	}

	return updateMask, nil
}
