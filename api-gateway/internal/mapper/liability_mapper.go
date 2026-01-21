package mapper

import (
	"time"
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"
)

func TimeToPtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func ToLiabilityDomain(p *pb.Liability) *domain.Liability {
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

	return &domain.Liability{
		ID:           p.Id,
		Name:         p.Name,
		Principal:    p.Principal,
		Creditor:     p.Creditor,
		InterestRate: p.InterestRate,
		Description:  p.Description,
		Type:         p.Type.String(),
		UserID:       p.CreatedBy,
		StartAt:      TimeToPtr(p.StartAt.AsTime()),
		EndAt:        TimeToPtr(p.EndAt.AsTime()),
		Files:        domainFiles,
		CreatedAt:    p.CreatedAt.AsTime(),
		UpdatedAt:    p.UpdatedAt.AsTime(),
	}
}

func ToLiabilityList(pbList []*pb.Liability) []domain.Liability {
	if pbList == nil {
		return []domain.Liability{}
	}

	entities := make([]domain.Liability, 0, len(pbList))

	for _, pbItem := range pbList {
		if item := ToLiabilityDomain(pbItem); item != nil {
			entities = append(entities, *item)
		}
	}

	return entities
}
