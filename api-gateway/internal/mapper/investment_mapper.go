package mapper

import (
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"
)

func ToInvestDomain(p *pb.Investment) *domain.Investment {
	if p == nil {
		return nil
	}

	var domainFiles []domain.FileInfo
	if len(p.Files) > 0 {
		domainFiles = make([]domain.FileInfo, len(p.Files))
		for i, f := range p.Files {
			domainFiles[i] = domain.FileInfo{
				ID:       f.Id,
				URL:      f.Url,
				FileType: f.FileType,
			}
		}
	}

	return &domain.Investment{
		ID:           p.Id,
		Name:         p.Name,
		Type:         p.Type.String(),
		Symbol:       p.Symbol,
		BrokerName:   p.BrokerName,
		Quantity:     p.Quantity,
		CostPerPrice: p.CostPrice,
		Amount:       p.Amount,
		Description:  p.Description,
		UserID:       p.UserId,
		Files:        domainFiles,
		CreatedAt:    p.CreatedAt.AsTime(),
		UpdatedAt:    p.UpdatedAt.AsTime(),
	}
}

func ToInvestList(pbList []*pb.Investment) []domain.Investment {
	if pbList == nil {
		return []domain.Investment{}
	}

	entities := make([]domain.Investment, 0, len(pbList))

	for _, pbItem := range pbList {
		if item := ToInvestDomain(pbItem); item != nil {
			entities = append(entities, *item)
		}
	}

	return entities
}
