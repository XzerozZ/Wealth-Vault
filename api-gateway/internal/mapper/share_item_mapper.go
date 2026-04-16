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
			Type:        item.Type,
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
			Type:     item.Type,
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

			case *pb.AssetPreview_Deleted:
				dtoItem.Type = v.Deleted.OriginalType
				dtoItem.AssetDetail = domain.DeletedDetail{
					ID:      v.Deleted.Id,
					Name:    v.Deleted.OriginalName,
					Message: v.Deleted.Message,
				}
			}
		}

		response = append(response, dtoItem)
	}

	return response
}

func MapAllFriendItemsToDomain(protoItems []*pb.SharedAssetPreview) []domain.FriendItemPreviewResponse {
	var response []domain.FriendItemPreviewResponse

	for _, item := range protoItems {
		dtoItem := domain.FriendItemPreviewResponse{
			ItemID: item.Id,
			Type:   item.Type,
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
					Name:           v.Insurance.Name,
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
					Amount:   v.Investment.Amount,
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

			case *pb.AssetPreview_Deleted:
				dtoItem.Type = v.Deleted.OriginalType
				dtoItem.AssetDetail = domain.DeletedDetail{
					ID:      v.Deleted.Id,
					Name:    v.Deleted.OriginalName,
					Message: v.Deleted.Message,
				}
			}
		}

		response = append(response, dtoItem)
	}

	return response
}

func ToSharedTargetsResponse(pbResp *pb.GetItemSharedTargetsResponse) domain.SharedTargetsResponse {
	resp := domain.SharedTargetsResponse{
		Groups:  make([]domain.SharedGroupTarget, 0),
		Friends: make([]domain.SharedFriendTarget, 0),
		Emails:  make([]domain.SharedEmailTarget, 0),
	}

	if pbResp == nil {
		return resp
	}

	for _, g := range pbResp.Groups {
		resp.Groups = append(resp.Groups, domain.SharedGroupTarget{
			GroupID:     g.GroupId,
			GroupName:   g.GroupName,
			GroupImage:  g.GroupImage,
			MemberCount: g.MemberCount,
			SharedAt:    g.SharedAt.AsTime(),
		})
	}

	for _, f := range pbResp.Friends {
		resp.Friends = append(resp.Friends, domain.SharedFriendTarget{
			FriendID:     f.FriendId,
			Username:     f.Username,
			ProfileImage: f.ProfileImage,
			SharedAt:     f.SharedAt.AsTime(),
		})
	}

	for _, e := range pbResp.Emails {
		resp.Emails = append(resp.Emails, domain.SharedEmailTarget{
			Email:    e.Email,
			SharedAt: e.SharedAt.AsTime(),
			IsSent:   e.IsSent,
		})
	}

	return resp
}
