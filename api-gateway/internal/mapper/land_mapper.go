package mapper

import (
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"
)

func ToLandDomain(p *pb.Land) *domain.Land {
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

	return &domain.Land{
		ID:          p.Id,
		Name:        p.Name,
		DeedNum:     p.DeedNum,
		Area:        p.Area,
		Amount:      p.Amount,
		Description: p.Description,
		Location:    *ToLocationDomain(p.Location),
		Ref:         ref,
		UserID:      p.UserId,
		Files:       domainFiles,
		CreatedAt:   p.CreatedAt.AsTime(),
		UpdatedAt:   p.UpdatedAt.AsTime(),
	}
}

func ToLandList(pbList []*pb.Land) []domain.Land {
	if pbList == nil {
		return []domain.Land{}
	}

	entities := make([]domain.Land, 0, len(pbList))

	for _, pbItem := range pbList {
		if item := ToLandDomain(pbItem); item != nil {
			entities = append(entities, *item)
		}
	}

	return entities
}
