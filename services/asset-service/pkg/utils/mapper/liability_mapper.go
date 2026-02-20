package mapper

import (
	"strings"
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils"
	"wealth-vault/asset-service/pkg/utils/helper"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var domainToProtoLiaType = map[string]pb.LiabilityType{
	"LOAN":    pb.LiabilityType_LIABILITY_TYPE_LOAN,
	"EXPENSE": pb.LiabilityType_LIABILITY_TYPE_EXPENSE,
}

func ToLiabilityDomain(req *pb.CreateLiabilityRequest, userID uuid.UUID) *domain.Liability {
	liType := domain.LiabilityTypeLoan
	if val, ok := helper.ProtoToDomainLiType[req.Type]; ok {
		liType = val
	}

	return &domain.Liability{
		UserID:       userID,
		Type:         liType,
		Name:         req.Name,
		Creditor:     req.Creditor,
		Principal:    req.Principal,
		InterestRate: req.InterestRate,
		Description:  req.Description,
		StartAt:      utils.ToTimePtr(req.StartAt),
		EndAt:        utils.ToTimePtr(req.EndAt),
		Files:        ToDomainFiles(req.NewFiles, userID),
	}
}

func ToLiabilityProto(d *domain.Liability) *pb.Liability {
	normalizedType := strings.ToUpper(strings.TrimSpace(string(d.Type)))
	inTypeEnum, ok := domainToProtoLiaType[normalizedType]
	if !ok {
		inTypeEnum = pb.LiabilityType_LIABILITY_TYPE_UNSPECIFIED
	}

	var startPb, endPb *timestamppb.Timestamp
	if d.StartAt != nil {
		startPb = timestamppb.New(*d.StartAt)
	}

	if d.EndAt != nil {
		endPb = timestamppb.New(*d.EndAt)
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
		StartAt:      startPb,
		EndAt:        endPb,
		CreatedAt:    timestamppb.New(d.CreatedAt),
		UpdatedAt:    timestamppb.New(d.UpdatedAt),
		Files:        ToPbFiles(d.Files),
		DeletedAt:    timestamppb.New(d.DeletedAt.Time),
	}

	return res
}

func ToLiabilityProtoSlice(liabilities []*domain.Liability) []*pb.Liability {
	if len(liabilities) == 0 {
		return []*pb.Liability{}
	}

	res := make([]*pb.Liability, len(liabilities))
	for i, lia := range liabilities {
		res[i] = ToLiabilityProto(lia)
	}
	return res
}
