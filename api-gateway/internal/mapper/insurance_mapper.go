package mapper

import (
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"
)

func ToInsuranceDomain(p *pb.Insurance) *domain.Insurance {
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

	return &domain.Insurance{
		ID:             p.Id,
		UserID:         p.UserId,
		PolicyNumber:   p.PolNum,
		Type:           p.Type.String(),
		CompanyName:    p.CompanyName,
		CoveragePeriod: p.CoveragePeriod,
		CoverageAmount: p.CoverageAmount,
		Description:    p.Description,
		ConDate:        TimeToPtr(p.ConDate.AsTime()),
		ExpDate:        TimeToPtr(p.ExpDate.AsTime()),
		Files:          domainFiles,
		CreatedAt:      p.CreatedAt.AsTime(),
		UpdatedAt:      p.UpdatedAt.AsTime(),
	}
}

func ToInsuranceList(pbList []*pb.Insurance) []domain.Insurance {
	if pbList == nil {
		return []domain.Insurance{}
	}

	entities := make([]domain.Insurance, 0, len(pbList))

	for _, pbItem := range pbList {
		if item := ToInsuranceDomain(pbItem); item != nil {
			entities = append(entities, *item)
		}
	}

	return entities
}
