package helper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

func ApplyUpdateLandFields(req *pb.UpdateLandRequest, land *domain.Land) error {
	paths := req.UpdateMask.GetPaths()
	if len(paths) == 0 {
		return nil
	}

	for _, path := range paths {
		switch path {
		case "name":
			if req.Land.Name != "" {
				land.Name = req.Land.Name
			}

		case "deed_num":
			if req.Land.DeedNum != "" {
				land.DeedNum = req.Land.DeedNum
			}

		case "area":
			if req.Land.Area != 0 {
				land.Area = req.Land.Area
			}

		case "amount":
			if req.Land.Amount != 0 {
				land.Amount = req.Land.Amount
			}

		case "description":
			if req.Land.Description != "" {
				land.Description = req.Land.Description
			}

		case "location.address":
			if req.Land.Location.Address != "" {
				land.Location.Address = req.Land.Location.Address
			}

		case "location.sub_district":
			if req.Land.Location.Subdistrict != "" {
				land.Location.Subdistrict = req.Land.Location.Subdistrict
			}

		case "location.district":
			if req.Land.Location.District != "" {
				land.Location.District = req.Land.Location.District
			}

		case "location.province":
			if req.Land.Location.Province != "" {
				land.Location.Province = req.Land.Location.Province
			}

		case "location.postal_code":
			if req.Land.Location.PostalCode != "" {
				land.Location.PostalCode = req.Land.Location.PostalCode
			}

		}
	}

	return nil
}
