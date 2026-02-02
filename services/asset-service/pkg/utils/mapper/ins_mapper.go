package mapper

import (
	"strings"
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var domainToProtoInsType = map[string]pb.InsuranceType{
	"LIFE":     pb.InsuranceType_INSURANCE_TYPE_LIFE,
	"HEALTH":   pb.InsuranceType_INSURANCE_TYPE_HEALTH,
	"ACCIDENT": pb.InsuranceType_INSURANCE_TYPE_ACCIDENT,
	"PROPERTY": pb.InsuranceType_INSURANCE_TYPE_PROPERTY,
	"VEHICLE":  pb.InsuranceType_INSURANCE_TYPE_VEHICLE,
}

func ToInsuranceProto(d *domain.Insurance) *pb.Insurance {
	normalizedType := strings.ToUpper(strings.TrimSpace(string(d.Type)))
	inTypeEnum, ok := domainToProtoInsType[normalizedType]
	if !ok {
		inTypeEnum = pb.InsuranceType_INSURANCE_TYPE_UNSPECIFIED
	}

	var condatePb, expdatePb *timestamppb.Timestamp
	if d.ConDate != nil {
		condatePb = timestamppb.New(*d.ConDate)
	}

	if d.ExpDate != nil {
		expdatePb = timestamppb.New(*d.ExpDate)
	}

	res := &pb.Insurance{
		Id:             d.ID.String(),
		UserId:         d.UserID.String(),
		Name:           d.Name,
		PolNum:         d.PolicyNumber,
		CompanyName:    d.CompanyName,
		Type:           inTypeEnum,
		CoveragePeriod: d.CoveragePeriod,
		CoverageAmount: d.CoverageAmount,
		ConDate:        condatePb,
		ExpDate:        expdatePb,
		Description:    d.Description,
		Files:          ToPbFiles(d.Files),
		CreatedAt:      timestamppb.New(d.CreatedAt),
		UpdatedAt:      timestamppb.New(d.UpdatedAt),
	}

	return res
}
