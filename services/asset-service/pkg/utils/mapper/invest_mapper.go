package mapper

import (
	"strings"
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
	"wealth-vault/asset-service/pkg/utils/helper"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var domainToProtoInvestType = map[string]pb.InvestmentType{
	"STOCK_TH":    pb.InvestmentType_INVEST_TYPE_STOCK_TH,
	"STOCK_US":    pb.InvestmentType_INVEST_TYPE_STOCK_US,
	"MUTUAL_FUND": pb.InvestmentType_INVEST_TYPE_MUTUAL_FUND,
	"BOND":        pb.InvestmentType_INVEST_TYPE_BOND,
	"CRYPTO":      pb.InvestmentType_INVEST_TYPE_CRYPTO,
	"GOLD":        pb.InvestmentType_INVEST_TYPE_GOLD,
}

func ToInvestmentDomain(req *pb.CreateInvestmentRequest, userID uuid.UUID) *domain.Investment {
	inType := domain.InvestTypeStockUS
	if val, ok := helper.ProtoToDomainInType[req.Type]; ok {
		inType = val
	}

	return &domain.Investment{
		UserID:       userID,
		Name:         req.Name,
		Symbol:       req.Symbol,
		Type:         inType,
		BrokerName:   req.BrokerName,
		Quantity:     req.Quantity,
		CostPerPrice: req.CostPrice,
		Description:  req.Description,
		Files:        ToDomainFiles(req.NewFiles, userID),
	}
}

func ToInvestProto(d *domain.Investment) *pb.Investment {
	normalizedType := strings.ToUpper(strings.TrimSpace(string(d.Type)))
	inTypeEnum, ok := domainToProtoInvestType[normalizedType]
	if !ok {
		inTypeEnum = pb.InvestmentType_INVEST_TYPE_UNSPECIFIED
	}

	res := &pb.Investment{
		Id:          d.ID.String(),
		UserId:      d.UserID.String(),
		Name:        d.Name,
		Symbol:      d.Symbol,
		Type:        inTypeEnum,
		BrokerName:  d.BrokerName,
		Quantity:    d.Quantity,
		CostPrice:   d.CostPerPrice,
		Amount:      d.Amount,
		Description: d.Description,
		Files:       ToPbFiles(d.Files),
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
		DeletedAt:   timestamppb.New(d.DeletedAt.Time),
	}

	return res
}

func ToInvestProtoSlice(invests []*domain.Investment) []*pb.Investment {
	if len(invests) == 0 {
		return []*pb.Investment{}
	}

	res := make([]*pb.Investment, len(invests))
	for i, in := range invests {
		res[i] = ToInvestProto(in)
	}
	return res
}
