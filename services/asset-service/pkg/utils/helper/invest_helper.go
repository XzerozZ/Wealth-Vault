package helper

import (
	"wealth-vault/asset-service/internal/domain"
	pb "wealth-vault/asset-service/pkg/pb/proto/asset"
)

var ProtoToDomainInType = map[pb.InvestmentType]domain.InvestmentType{
	pb.InvestmentType_INVEST_TYPE_STOCK_TH:    domain.InvestTypeStockTH,
	pb.InvestmentType_INVEST_TYPE_STOCK_US:    domain.InvestTypeStockUS,
	pb.InvestmentType_INVEST_TYPE_MUTUAL_FUND: domain.InvestTypeMutualFund,
	pb.InvestmentType_INVEST_TYPE_BOND:        domain.InvestTypeBond,
	pb.InvestmentType_INVEST_TYPE_CRYPTO:      domain.InvestTypeCrypto,
	pb.InvestmentType_INVEST_TYPE_GOLD:        domain.InvestTypeGold,
}

func ApplyUpdateInFields(req *pb.UpdateInvestmentRequest, in *domain.Investment) error {
	paths := req.UpdateMask.GetPaths()
	if len(paths) == 0 {
		return nil
	}

	for _, path := range paths {
		switch path {
		case "name":
			if req.Invest.Name != "" {
				in.Name = req.Invest.Name
			}

		case "symbol":
			if req.Invest.Symbol != "" {
				in.Symbol = req.Invest.Symbol
			}

		case "broker_name":
			if req.Invest.BrokerName != "" {
				in.BrokerName = req.Invest.BrokerName
			}

		case "quantity":
			if req.Invest.Quantity != 0 {
				in.Quantity = req.Invest.Quantity
			}

		case "cost_per_price":
			if req.Invest.CostPrice != 0 {
				in.CostPerPrice = req.Invest.CostPrice
			}

		case "amount":
			if req.Invest.Amount != 0 {
				in.Amount = req.Invest.Amount
			}

		case "description":
			if req.Invest.Description != "" {
				in.Description = req.Invest.Description
			}

		case "type":
			if req.Invest.Type != pb.InvestmentType_INVEST_TYPE_UNSPECIFIED {
				if val, ok := ProtoToDomainInType[req.Invest.Type]; ok {
					in.Type = val
				}
			}
		}
	}

	return nil
}
