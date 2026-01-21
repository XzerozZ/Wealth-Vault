package helper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

func ApplyUpdateFields(req *pb.UpdateLiabilityRequest, lia *domain.Liability) ([]string, error) {
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
		lia.Name = req.Name
		updateMask = append(updateMask, "Name")
	}

	if has("creditor") {
		lia.Creditor = req.Creditor
		updateMask = append(updateMask, "Creditor")
	}

	if has("principal") {
		lia.Principal = req.Principal
		updateMask = append(updateMask, "Principal")
	}

	if has("interest_rate") {
		lia.InterestRate = req.InterestRate
		updateMask = append(updateMask, "InterestRate")
	}

	if has("description") {
		lia.Description = req.Description
		updateMask = append(updateMask, "Description")
	}

	if has("started_at") && req.StartAt != nil {
		t := req.StartAt.AsTime()
		lia.StartAt = &t
		updateMask = append(updateMask, "StartAt")
	}

	if has("ended_at") && req.EndAt != nil {
		t := req.EndAt.AsTime()
		lia.EndAt = &t
		updateMask = append(updateMask, "EndAt")
	}

	return updateMask, nil
}
