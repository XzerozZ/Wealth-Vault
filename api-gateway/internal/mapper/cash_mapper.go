package mapper

import (
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"
)

func ToCashDomain(p *pb.Cash) *domain.Cash {
	if p == nil {
		return nil
	}

	return &domain.Cash{
		ID:          p.Id,
		Name:        p.Name,
		Amount:      p.Amount,
		Description: p.Description,
		CreatedBy:   p.CreatedBy,
		CreatedAt:   p.CreatedAt.AsTime(),
		UpdatedAt:   p.UpdatedAt.AsTime(),
	}
}

func ToCashList(pbList []*pb.Cash) []domain.Cash {
	if pbList == nil {
		return []domain.Cash{}
	}

	entities := make([]domain.Cash, 0, len(pbList))

	for _, pbItem := range pbList {
		entities = append(entities, *ToCashDomain(pbItem)) // dereference ถ้า ToCashEntity คืนค่าเป็น pointer
	}

	return entities
}
