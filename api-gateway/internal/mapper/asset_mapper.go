package mapper

import (
	"strings"
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"
)

func SafeMapEnum(valueMap map[string]int32, key string, prefix string) int32 {
	key = strings.ToUpper(strings.TrimSpace(key))
	if v, ok := valueMap[key]; ok {
		return v
	}

	if v, ok := valueMap[prefix+key]; ok {
		return v
	}

	return 0
}

func ToAssetDomain(p *pb.AssetSummary) *domain.AssetSummary {
	if p == nil {
		return nil
	}

	return &domain.AssetSummary{
		ID:        p.Id,
		Name:      p.Name,
		Amount:    p.Value,
		Type:      p.Type,
		CreatedAt: p.CreatedAt.AsTime(),
	}
}

func ToAssetList(pbList []*pb.AssetSummary) []domain.AssetSummary {
	if pbList == nil {
		return []domain.AssetSummary{}
	}

	entities := make([]domain.AssetSummary, 0, len(pbList))

	for _, pbItem := range pbList {
		if item := ToAssetDomain(pbItem); item != nil {
			entities = append(entities, *item)
		}
	}

	return entities
}
