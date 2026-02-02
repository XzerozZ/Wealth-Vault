package mapper

import (
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"
)

func ToBuildingDomain(p *pb.Building) *domain.Building {
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

	var ref []domain.RefInfo
	if len(p.Ref) > 0 {
		ref = make([]domain.RefInfo, len(p.Ref))
		for i, f := range p.Ref {
			ref[i] = domain.RefInfo{
				ID:   f.Id,
				Name: f.Name,
			}
		}
	}

	var ins []domain.InsInfo
	if len(p.Ins) > 0 {
		ins = make([]domain.InsInfo, len(p.Ins))
		for i, f := range p.Ins {
			ins[i] = domain.InsInfo{
				ID:   f.Id,
				Name: f.Name,
			}
		}
	}

	return &domain.Building{
		ID:          p.Id,
		Name:        p.Name,
		Type:        p.Type.String(),
		Area:        p.Area,
		Amount:      p.Amount,
		Description: p.Description,
		Location:    *ToLocationDomain(p.Location),
		Ref:         ref,
		Ins:         ins,
		UserID:      p.UserId,
		Files:       domainFiles,
		CreatedAt:   p.CreatedAt.AsTime(),
		UpdatedAt:   p.UpdatedAt.AsTime(),
	}
}

func ToBuildingList(pbList []*pb.Building) []domain.Building {
	if pbList == nil {
		return []domain.Building{}
	}

	entities := make([]domain.Building, 0, len(pbList))

	for _, pbItem := range pbList {
		if item := ToBuildingDomain(pbItem); item != nil {
			entities = append(entities, *item)
		}
	}

	return entities
}
