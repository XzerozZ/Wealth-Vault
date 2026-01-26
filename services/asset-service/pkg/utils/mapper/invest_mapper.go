package mapper

import (
	"strings"
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"

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
		Description: d.Description,
		Files:       ToPbFiles(d.Files),
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
	}

	return res
}
