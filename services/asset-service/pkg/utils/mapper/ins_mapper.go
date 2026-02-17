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

var domainToProtoInsType = map[string]pb.InsuranceType{
	"LIFE":     pb.InsuranceType_INSURANCE_TYPE_LIFE,
	"HEALTH":   pb.InsuranceType_INSURANCE_TYPE_HEALTH,
	"ACCIDENT": pb.InsuranceType_INSURANCE_TYPE_ACCIDENT,
	"PROPERTY": pb.InsuranceType_INSURANCE_TYPE_PROPERTY,
	"VEHICLE":  pb.InsuranceType_INSURANCE_TYPE_VEHICLE,
}

func MapToInsurances(ids []uuid.UUID) []domain.Insurance {
	res := make([]domain.Insurance, len(ids))
	for i, id := range ids {
		res[i] = domain.Insurance{
			ID: id,
		}
	}
	return res
}

func ToInsuranceDomain(req *pb.CreateInsuranceRequest, userID uuid.UUID) *domain.Insurance {
	inType := domain.InsuranceTypeLife
	if val, ok := helper.ProtoToDomainInsType[req.Type]; ok {
		inType = val
	}

	return &domain.Insurance{
		UserID:         userID,
		Type:           inType,
		Name:           req.Name,
		PolicyNumber:   req.PolNum,
		CompanyName:    req.CompanyName,
		CoveragePeriod: req.CoveragePeriod,
		CoverageAmount: req.CoverageAmount,
		ConDate:        utils.ToTimePtr(req.ConDate),
		ExpDate:        utils.ToTimePtr(req.ExpDate),
		Description:    req.Description,
		Files:          ToDomainFiles(req.NewFiles, userID),
	}
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
		DeletedAt:      timestamppb.New(d.DeletedAt.Time),
	}

	return res
}

func ToInsuranceProtoSlice(insurances []*domain.Insurance) []*pb.Insurance {
	if len(insurances) == 0 {
		return []*pb.Insurance{}
	}

	res := make([]*pb.Insurance, len(insurances))
	for i, in := range insurances {
		res[i] = ToInsuranceProto(in)
	}
	return res
}
