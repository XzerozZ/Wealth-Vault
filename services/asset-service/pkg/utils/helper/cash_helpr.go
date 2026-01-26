package helper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

func ApplyUpdateCashFields(req *pb.UpdateCashRequest, cash *domain.Cash) error {
	paths := req.UpdateMask.GetPaths()
	if len(paths) == 0 {
		return nil
	}

	for _, path := range paths {
		switch path {
		case "name":
			if req.Cash.Name != "" {
				cash.Name = req.Cash.Name
			}

		case "amount":
			if req.Cash.Amount != 0 {
				cash.Amount = req.Cash.Amount
			}

		case "description":
			if req.Cash.Description != "" {
				cash.Description = req.Cash.Description
			}
		}
	}

	return nil
}
