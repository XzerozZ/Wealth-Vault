package mapper

import (
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/user"
)

func MapGroupItemsToDomain(protoItems []*pb.GroupItemDetail) []domain.GroupItemResponse {
	var response []domain.GroupItemResponse

	for _, item := range protoItems {
		dtoItem := domain.GroupItemResponse{
			GroupItemID: item.GroupItemId,
			SharedBy:    item.SharedBy,
			SharedAt:    item.SharedAt.AsTime(),
		}

		if item.AssetDetail != nil {
			switch v := item.AssetDetail.Asset.(type) {

			case *pb.AssetPreview_Building:
				dtoItem.Type = "building"
				dtoItem.AssetDetail = domain.BuildingDetail{
					ID:           v.Building.Id,
					Name:         v.Building.Name,
					Amount:       v.Building.Amount,
					LocationText: v.Building.LocationText,
					TypeName:     v.Building.TypeName,
				}

			case *pb.AssetPreview_Account:
				dtoItem.Type = "account"
				dtoItem.AssetDetail = domain.AccountDetail{
					ID:            v.Account.Id,
					Name:          v.Account.Name,
					BankName:      v.Account.BankName,
					AccountNumber: v.Account.AccountNumber,
					Amount:        v.Account.Amount,
				}

			case *pb.AssetPreview_Land:
				dtoItem.Type = "land"
				dtoItem.AssetDetail = domain.LandDetail{
					ID:           v.Land.Id,
					Name:         v.Land.Name,
					DeedNum:      v.Land.DeedNum,
					Area:         v.Land.Area,
					Amount:       v.Land.Amount,
					LocationText: v.Land.LocationText,
				}

			case *pb.AssetPreview_Cash:
				dtoItem.Type = "cash"
				dtoItem.AssetDetail = domain.CashDetail{
					ID:     v.Cash.Id,
					Name:   v.Cash.Name,
					Amount: v.Cash.Amount,
				}

			case *pb.AssetPreview_Insurance:
				dtoItem.Type = "insurance"
				dtoItem.AssetDetail = domain.InsuranceDetail{
					ID:             v.Insurance.Id,
					CompanyName:    v.Insurance.CompanyName,
					TypeName:       v.Insurance.TypeName,
					PolNum:         v.Insurance.PolNum,
					CoverageAmount: v.Insurance.CoverageAmount,
					ExpDateText:    v.Insurance.ExpDateText,
				}

			case *pb.AssetPreview_Investment:
				dtoItem.Type = "investment"
				dtoItem.AssetDetail = domain.InvestmentDetail{
					ID:       v.Investment.Id,
					Name:     v.Investment.Name,
					Symbol:   v.Investment.Symbol,
					TypeName: v.Investment.TypeName,
				}

			case *pb.AssetPreview_Liability:
				dtoItem.Type = "liability"
				dtoItem.AssetDetail = domain.LiabilityDetail{
					ID:        v.Liability.Id,
					Name:      v.Liability.Name,
					Creditor:  v.Liability.Creditor,
					Principal: v.Liability.Principal,
					TypeName:  v.Liability.TypeName,
				}
			}
		}

		response = append(response, dtoItem)
	}

	return response
}

func MapFriendItemsToDomain(protoItems []*pb.FriendItemDetail) []domain.FriendItemResponse {
	var response []domain.FriendItemResponse

	for _, item := range protoItems {
		dtoItem := domain.FriendItemResponse{
			ItemID:   item.FriendItemId,
			SharedBy: item.SharedBy,
			SharedAt: item.SharedAt.AsTime(),
		}

		if item.AssetDetail != nil {
			switch v := item.AssetDetail.Asset.(type) {

			case *pb.AssetPreview_Building:
				dtoItem.Type = "building"
				dtoItem.AssetDetail = domain.BuildingDetail{
					ID:           v.Building.Id,
					Name:         v.Building.Name,
					Amount:       v.Building.Amount,
					LocationText: v.Building.LocationText,
					TypeName:     v.Building.TypeName,
				}

			case *pb.AssetPreview_Account:
				dtoItem.Type = "account"
				dtoItem.AssetDetail = domain.AccountDetail{
					ID:            v.Account.Id,
					Name:          v.Account.Name,
					BankName:      v.Account.BankName,
					AccountNumber: v.Account.AccountNumber,
					Amount:        v.Account.Amount,
				}

			case *pb.AssetPreview_Land:
				dtoItem.Type = "land"
				dtoItem.AssetDetail = domain.LandDetail{
					ID:           v.Land.Id,
					Name:         v.Land.Name,
					DeedNum:      v.Land.DeedNum,
					Area:         v.Land.Area,
					Amount:       v.Land.Amount,
					LocationText: v.Land.LocationText,
				}

			case *pb.AssetPreview_Cash:
				dtoItem.Type = "cash"
				dtoItem.AssetDetail = domain.CashDetail{
					ID:     v.Cash.Id,
					Name:   v.Cash.Name,
					Amount: v.Cash.Amount,
				}

			case *pb.AssetPreview_Insurance:
				dtoItem.Type = "insurance"
				dtoItem.AssetDetail = domain.InsuranceDetail{
					ID:             v.Insurance.Id,
					CompanyName:    v.Insurance.CompanyName,
					TypeName:       v.Insurance.TypeName,
					PolNum:         v.Insurance.PolNum,
					CoverageAmount: v.Insurance.CoverageAmount,
					ExpDateText:    v.Insurance.ExpDateText,
				}

			case *pb.AssetPreview_Investment:
				dtoItem.Type = "investment"
				dtoItem.AssetDetail = domain.InvestmentDetail{
					ID:       v.Investment.Id,
					Name:     v.Investment.Name,
					Symbol:   v.Investment.Symbol,
					TypeName: v.Investment.TypeName,
				}

			case *pb.AssetPreview_Liability:
				dtoItem.Type = "liability"
				dtoItem.AssetDetail = domain.LiabilityDetail{
					ID:        v.Liability.Id,
					Name:      v.Liability.Name,
					Creditor:  v.Liability.Creditor,
					Principal: v.Liability.Principal,
					TypeName:  v.Liability.TypeName,
				}
			}
		}

		response = append(response, dtoItem)
	}

	return response
}
