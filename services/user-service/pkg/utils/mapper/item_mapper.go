package mapper

import (
	"fmt"
	"time"
	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils"
)

// const GhostDuration = 7 * 24 * time.Hour
const GhostDuration = 2 * time.Minute

func createGhostPreview(id string, name string, itemtype string) *pb.AssetPreview {
	return &pb.AssetPreview{
		Asset: &pb.AssetPreview_Deleted{
			Deleted: &pb.DeletedPreview{
				Id:           id,
				OriginalName: name,
				OriginalType: itemtype,
				Message:      "รายการนี้ถูกลบไปแล้ว : ไม่สามารถแสดงรายละเอียดได้",
			},
		},
	}
}

func MapBuildingToPreview(b *assetPb.Building) *pb.AssetPreview {
	if b == nil {
		return nil
	}

	if b.DeletedAt != nil && b.DeletedAt.IsValid() && b.DeletedAt.Seconds > 0 {
		t := b.DeletedAt.AsTime()
		if t.Unix() <= 0 {
			goto ProcessActive
		}

		if time.Since(t) > GhostDuration {
			return nil
		}

		return createGhostPreview(b.Id, b.Name, "building")
	}

ProcessActive:
	locationText := "ไม่ระบุตำแหน่ง"
	if b.Location != nil {
		locationText = fmt.Sprintf("%s, %s", b.Location.District, b.Location.Province)
	}

	imageURL := ""
	if len(b.Files) > 0 && b.Files[0] != nil {
		imageURL = b.Files[0].Url
	}

	return &pb.AssetPreview{
		Asset: &pb.AssetPreview_Building{
			Building: &pb.BuildingPreview{
				Id:           b.Id,
				Name:         b.Name,
				Amount:       b.Amount,
				LocationText: locationText,
				TypeName:     b.Type.String(),
				Image:        imageURL,
			},
		},
	}
}

func MapLandToPreview(l *assetPb.Land) *pb.AssetPreview {
	if l == nil {
		return nil
	}

	if l.DeletedAt != nil && l.DeletedAt.IsValid() && l.DeletedAt.Seconds > 0 {
		t := l.DeletedAt.AsTime()
		if t.Unix() <= 0 {
			goto ProcessActive
		}

		if time.Since(t) > GhostDuration {
			return nil
		}
		return createGhostPreview(l.Id, l.Name, "land")
	}

ProcessActive:
	locationText := "ไม่ระบุตำแหน่ง"
	if l.Location != nil {
		locationText = fmt.Sprintf("%s, %s", l.Location.District, l.Location.Province)
	}

	imageURL := ""
	if len(l.Files) > 0 && l.Files[0] != nil {
		imageURL = l.Files[0].Url
	}

	return &pb.AssetPreview{
		Asset: &pb.AssetPreview_Land{
			Land: &pb.LandPreview{
				Id:           l.Id,
				Name:         l.Name,
				DeedNum:      l.DeedNum,
				Area:         l.Area,
				Amount:       l.Amount,
				LocationText: locationText,
				Image:        imageURL,
			},
		},
	}
}

func MapAccountToPreview(a *assetPb.Account) *pb.AssetPreview {
	if a == nil {
		return nil
	}

	fmt.Printf("DEBUG MAPPER: ID=%s, DeletedAt=%v\n", a.Id, a.DeletedAt)

	if a.DeletedAt != nil && a.DeletedAt.IsValid() && a.DeletedAt.Seconds > 0 {
		t := a.DeletedAt.AsTime()
		if t.Unix() <= 0 {
			goto ProcessActive
		}

		if time.Since(t) > GhostDuration {
			return nil
		}
		return createGhostPreview(a.Id, a.Name, "account")
	}

ProcessActive:
	imageURL := ""
	if len(a.Files) > 0 && a.Files[0] != nil {
		imageURL = a.Files[0].Url
	}

	return &pb.AssetPreview{
		Asset: &pb.AssetPreview_Account{
			Account: &pb.AccountPreview{
				Id:            a.Id,
				Name:          a.Name,
				BankName:      a.BankName,
				Amount:        a.Amount,
				AccountNumber: utils.MaskBankAccount(a.BankAcc),
				Image:         imageURL,
			},
		},
	}
}

func MapCashToPreview(c *assetPb.Cash) *pb.AssetPreview {
	if c == nil {
		return nil
	}

	if c.DeletedAt != nil && c.DeletedAt.IsValid() && c.DeletedAt.Seconds > 0 {
		t := c.DeletedAt.AsTime()
		if t.Unix() <= 0 {
			goto ProcessActive
		}

		if time.Since(t) > GhostDuration {
			return nil
		}
		return createGhostPreview(c.Id, c.Name, "cash")
	}

ProcessActive:
	imageURL := ""
	if len(c.Files) > 0 && c.Files[0] != nil {
		imageURL = c.Files[0].Url
	}

	return &pb.AssetPreview{
		Asset: &pb.AssetPreview_Cash{
			Cash: &pb.CashPreview{
				Id:     c.Id,
				Name:   c.Name,
				Amount: c.Amount,
				Image:  imageURL,
			},
		},
	}
}

func MapInsuranceToPreview(i *assetPb.Insurance) *pb.AssetPreview {
	if i == nil {
		return nil
	}

	if i.DeletedAt != nil && i.DeletedAt.IsValid() && i.DeletedAt.Seconds > 0 {
		t := i.DeletedAt.AsTime()
		if t.Unix() <= 0 {
			goto ProcessActive
		}

		if time.Since(t) > GhostDuration {
			return nil
		}
		return createGhostPreview(i.Id, i.Name, "insurance")
	}

ProcessActive:
	expDateStr := "ไม่ระบุ"
	if i.ExpDate != nil {
		expDateStr = i.ExpDate.AsTime().String()
	}

	imageURL := ""
	if len(i.Files) > 0 && i.Files[0] != nil {
		imageURL = i.Files[0].Url
	}

	return &pb.AssetPreview{
		Asset: &pb.AssetPreview_Insurance{
			Insurance: &pb.InsurancePreview{
				Id:             i.Id,
				TypeName:       i.Type.String(),
				Name:           i.Name,
				CompanyName:    i.CompanyName,
				PolNum:         i.PolNum,
				CoverageAmount: i.CoverageAmount,
				ExpDateText:    expDateStr,
				Image:          imageURL,
			},
		},
	}
}

func MapInvestmentToPreview(inv *assetPb.Investment) *pb.AssetPreview {
	if inv == nil {
		return nil
	}

	if inv.DeletedAt != nil && inv.DeletedAt.IsValid() && inv.DeletedAt.Seconds > 0 {
		t := inv.DeletedAt.AsTime()
		if t.Unix() <= 0 {
			goto ProcessActive
		}

		if time.Since(t) > GhostDuration {
			return nil
		}
		return createGhostPreview(inv.Id, inv.Name, "investment")
	}

ProcessActive:
	imageURL := ""
	if len(inv.Files) > 0 && inv.Files[0] != nil {
		imageURL = inv.Files[0].Url
	}

	return &pb.AssetPreview{
		Asset: &pb.AssetPreview_Investment{
			Investment: &pb.InvestmentPreview{
				Id:       inv.Id,
				Name:     inv.Name,
				TypeName: inv.Type.String(),
				Symbol:   inv.Symbol,
				Amount:   inv.Amount,
				Image:    imageURL,
			},
		},
	}
}

func MapLiabilityToPreview(l *assetPb.Liability) *pb.AssetPreview {
	if l == nil {
		return nil
	}

	if l.DeletedAt != nil && l.DeletedAt.IsValid() && l.DeletedAt.Seconds > 0 {
		t := l.DeletedAt.AsTime()
		if t.Unix() <= 0 {
			goto ProcessActive
		}

		if time.Since(t) > GhostDuration {
			return nil
		}
		return createGhostPreview(l.Id, l.Name, "liability")
	}

ProcessActive:
	imageURL := ""
	if len(l.Files) > 0 && l.Files[0] != nil {
		imageURL = l.Files[0].Url
	}

	return &pb.AssetPreview{
		Asset: &pb.AssetPreview_Liability{
			Liability: &pb.LiabilityPreview{
				Id:        l.Id,
				Name:      l.Name,
				TypeName:  l.Type.String(),
				Creditor:  l.Creditor,
				Principal: l.Principal,
				Image:     imageURL,
			},
		},
	}
}
