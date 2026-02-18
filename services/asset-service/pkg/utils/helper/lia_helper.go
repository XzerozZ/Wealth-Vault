package helper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

var ProtoToDomainLiType = map[pb.LiabilityType]domain.LiabilityType{
	pb.LiabilityType_LIABILITY_TYPE_LOAN:    domain.LiabilityTypeLoan,
	pb.LiabilityType_LIABILITY_TYPE_EXPENSE: domain.LiabilityTypeExpense,
}

func ApplyUpdateFields(req *pb.UpdateLiabilityRequest, lia *domain.Liability) error {
	paths := req.UpdateMask.GetPaths()
	if len(paths) == 0 {
		return nil
	}

	for _, path := range paths {
		switch path {
		case "name":
			if req.Liability.Name != "" {
				lia.Name = req.Liability.Name
			}

		case "creditor":
			if req.Liability.Creditor != "" {
				lia.Creditor = req.Liability.Creditor
			}

		case "principal":
			if req.Liability.Principal != 0 {
				lia.Principal = req.Liability.Principal
			}

		case "interest_rate":
			if req.Liability.InterestRate != 0 {
				lia.InterestRate = req.Liability.InterestRate
			}

		case "started_at":
			if req.Liability.StartAt != nil {
				t := req.Liability.StartAt.AsTime()
				lia.StartAt = &t
			}

		case "ended_at":
			if req.Liability.EndAt != nil {
				t := req.Liability.EndAt.AsTime()
				lia.EndAt = &t
			}
		case "description":
			if req.Liability.Description != "" {
				lia.Description = req.Liability.Description
			}

		case "type":
			if req.Liability.Type != pb.LiabilityType_LIABILITY_TYPE_UNSPECIFIED {
				if val, ok := ProtoToDomainLiType[req.Liability.Type]; ok {
					lia.Type = val
				}
			}
		}

	}

	return nil
}
