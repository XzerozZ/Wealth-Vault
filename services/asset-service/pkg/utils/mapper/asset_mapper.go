package mapper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToAssetSummaryProto(d domain.AssetSummary) *pb.AssetSummary {
	imageURL := ""
	if len(d.Files) > 0 {
		if d.Files[0].Link != "" {
			imageURL = d.Files[0].Link
		}
	}

	return &pb.AssetSummary{
		Id:        d.ID.String(),
		Type:      d.Type,
		Name:      d.Name,
		Value:     d.Value,
		Image:     imageURL,
		CreatedAt: timestamppb.New(d.CreatedAt),
	}
}

func ToAssetSummaryProtoList(domains []domain.AssetSummary) []*pb.AssetSummary {
	var protos []*pb.AssetSummary
	for _, d := range domains {
		protos = append(protos, ToAssetSummaryProto(d))
	}
	return protos
}
