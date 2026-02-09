package mapper

import (
	"strings"
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var domainToProtoLiaType = map[string]pb.LiabilityType{
	"LOAN":    pb.LiabilityType_LIABILITY_TYPE_LOAN,
	"EXPENSE": pb.LiabilityType_LIABILITY_TYPE_EXPENSE,
}

func ToLiabilityProto(d *domain.Liability) *pb.Liability {
	normalizedType := strings.ToUpper(strings.TrimSpace(string(d.Type)))
	inTypeEnum, ok := domainToProtoLiaType[normalizedType]
	if !ok {
		inTypeEnum = pb.LiabilityType_LIABILITY_TYPE_UNSPECIFIED
	}

	res := &pb.Liability{
		Id:           d.ID.String(),
		Name:         d.Name,
		Principal:    d.Principal,
		InterestRate: d.InterestRate,
		Creditor:     d.Creditor,
		Type:         inTypeEnum,
		Description:  d.Description,
		CreatedBy:    d.UserID.String(),
		StartAt:      timestamppb.New(*d.StartAt),
		EndAt:        timestamppb.New(*d.EndAt),
		CreatedAt:    timestamppb.New(d.CreatedAt),
		UpdatedAt:    timestamppb.New(d.UpdatedAt),
		Files:        ToPbFiles(d.Files),
		DeletedAt:    timestamppb.New(d.DeletedAt.Time),
	}

	return res
}
