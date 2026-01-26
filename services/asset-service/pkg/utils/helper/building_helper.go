package helper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

var ProtoToDomainBuildingType = map[pb.BuildingType]domain.BuildingType{
	pb.BuildingType_BUILDING_TYPE_CONDO:      domain.BuildingTypeCondo,
	pb.BuildingType_BUILDING_TYPE_HOUSE:      domain.BuildingTypeHouse,
	pb.BuildingType_BUILDING_TYPE_TOWNHOME:   domain.BuildingTypeTownHome,
	pb.BuildingType_BUILDING_TYPE_COMMERCIAL: domain.BuildingTypeCommercial,
}

func ApplyUpdateBuildingFields(req *pb.UpdateBuildingRequest, build *domain.Building) error {
	paths := req.UpdateMask.GetPaths()
	if len(paths) == 0 {
		return nil
	}

	for _, path := range paths {
		switch path {
		case "name":
			if req.Building.Name != "" {
				build.Name = req.Building.Name
			}

		case "area":
			if req.Building.Area != 0 {
				build.Area = req.Building.Area
			}

		case "amount":
			if req.Building.Amount != 0 {
				build.Amount = req.Building.Amount
			}

		case "description":
			if req.Building.Description != "" {
				build.Description = req.Building.Description
			}

		case "location.address":
			if req.Building.Location.Address != "" {
				build.Location.Address = req.Building.Location.Address
			}

		case "location.sub_district":
			if req.Building.Location.Subdistrict != "" {
				build.Location.Subdistrict = req.Building.Location.Subdistrict
			}

		case "location.district":
			if req.Building.Location.District != "" {
				build.Location.District = req.Building.Location.District
			}

		case "location.province":
			if req.Building.Location.Province != "" {
				build.Location.Province = req.Building.Location.Province
			}

		case "location.postal_code":
			if req.Building.Location.PostalCode != "" {
				build.Location.PostalCode = req.Building.Location.PostalCode
			}

		case "type":
			if req.Building.Type != pb.BuildingType_BUILDING_TYPE_UNSPECIFIED {
				if val, ok := ProtoToDomainBuildingType[req.Building.Type]; ok {
					build.Type = val
				}
			}
		}
	}

	return nil
}
