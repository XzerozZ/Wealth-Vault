package mapper

import (
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"
)

func ToAccountDomain(p *pb.Account) *domain.Account {
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

	return &domain.Account{
		ID:          p.Id,
		Name:        p.Name,
		BankName:    p.BankName,
		BankAccount: p.BankAcc,
		Amount:      p.Amount,
		Description: p.Description,
		Type:        p.Type.String(),
		UserID:      p.UserId,
		Files:       domainFiles,
		CreatedAt:   p.CreatedAt.AsTime(),
		UpdatedAt:   p.UpdatedAt.AsTime(),
	}
}

func ToAccountList(pbList []*pb.Account) []domain.Account {
	if pbList == nil {
		return []domain.Account{}
	}

	entities := make([]domain.Account, 0, len(pbList))

	for _, pbItem := range pbList {
		if item := ToAccountDomain(pbItem); item != nil {
			entities = append(entities, *item)
		}
	}

	return entities
}
