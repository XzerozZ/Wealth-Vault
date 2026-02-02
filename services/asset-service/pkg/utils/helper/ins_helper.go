package helper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

var ProtoToDomainInsType = map[pb.InsuranceType]domain.InsuranceType{
	pb.InsuranceType_INSURANCE_TYPE_LIFE:     domain.InsuranceTypeLife,
	pb.InsuranceType_INSURANCE_TYPE_HEALTH:   domain.InsuranceTypeHealth,
	pb.InsuranceType_INSURANCE_TYPE_ACCIDENT: domain.InsuranceTypeAccident,
	pb.InsuranceType_INSURANCE_TYPE_PROPERTY: domain.InsuranceTypeProperty,
	pb.InsuranceType_INSURANCE_TYPE_VEHICLE:  domain.InsuranceTypeVehicle,
}

func ApplyUpdateInsuranceFields(req *pb.UpdateInsuranceRequest, in *domain.Insurance) error {
	paths := req.UpdateMask.GetPaths()
	if len(paths) == 0 {
		return nil
	}

	for _, path := range paths {
		switch path {
		case "name":
			if req.Insurance.Name != "" {
				in.Name = req.Insurance.Name
			}

		case "policy_number":
			if req.Insurance.PolNum != "" {
				in.PolicyNumber = req.Insurance.PolNum
			}

		case "company_name":
			if req.Insurance.CompanyName != "" {
				in.CompanyName = req.Insurance.CompanyName
			}

		case "coverage_period":
			if req.Insurance.CoveragePeriod != 0 {
				in.CoveragePeriod = req.Insurance.CoveragePeriod
			}

		case "coverage_amount":
			if req.Insurance.CoverageAmount != 0 {
				in.CoverageAmount = req.Insurance.CoverageAmount
			}

		case "con_date":
			if req.Insurance.ConDate != nil {
				t := req.Insurance.ConDate.AsTime()
				in.ConDate = &t
			}

		case "exp_date":
			if req.Insurance.ExpDate != nil {
				t := req.Insurance.ExpDate.AsTime()
				in.ExpDate = &t
			}

		case "description":
			if req.Insurance.Description != "" {
				in.Description = req.Insurance.Description
			}

		case "type":
			if req.Insurance.Type != pb.InsuranceType_INSURANCE_TYPE_UNSPECIFIED {
				if val, ok := ProtoToDomainInsType[req.Insurance.Type]; ok {
					in.Type = val
				}
			}
		}
	}

	return nil
}
