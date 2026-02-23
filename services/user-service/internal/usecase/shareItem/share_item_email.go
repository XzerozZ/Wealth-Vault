package usecase

import (
	"context"
	"fmt"
	"time"
	"wealth-vault/user-service/internal/domain"
	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"
	"wealth-vault/user-service/pkg/utils"

	"github.com/google/uuid"
)

type AssetInfo struct {
	Name    string
	Image   string
	Details map[string]string
}

func (u *ShareItemUsecase) SendEmailInvitations(items []domain.EmailItem) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idsByType := make(map[string][]string)
	for _, item := range items {
		idsByType[item.EntityType] = append(idsByType[item.EntityType], item.EntityID.String())
	}

	infoMap := make(map[string]AssetInfo)

	if len(idsByType[AssetTypeBuilding]) > 0 {
		if res, _ := u.assetClient.GetBatchBuilding(ctx, &assetPb.GetBatchIdsRequest{Ids: idsByType[AssetTypeBuilding]}); res != nil {
			for _, b := range res.Building {
				locationStr := "ไม่ระบุตำแหน่ง"
				if b.Location != nil {
					locationStr = fmt.Sprintf("%s, %s", b.Location.District, b.Location.Province)
				}

				infoMap[b.Id] = AssetInfo{Name: b.Name, Details: map[string]string{
					"ชื่อ":    b.Name,
					"มูลค่า":  fmt.Sprintf("฿%.2f", b.Amount),
					"ที่ตั้ง": locationStr,
					"ประเภท":  b.Type.String(),
				}}
			}
		}
	}

	if len(idsByType[AssetTypeLand]) > 0 {
		if res, _ := u.assetClient.GetBatchLand(ctx, &assetPb.GetBatchIdsRequest{Ids: idsByType[AssetTypeLand]}); res != nil {
			for _, l := range res.Land {
				locationStr := "ไม่ระบุตำแหน่ง"
				if l.Location != nil {
					locationStr = fmt.Sprintf("%s, %s", l.Location.District, l.Location.Province)
				}

				infoMap[l.Id] = AssetInfo{Name: l.Name, Details: map[string]string{
					"ชื่อ":     l.Name,
					"เลขโฉนด":  l.DeedNum,
					"เนื้อที่": fmt.Sprintf("%.2f", l.Area),
					"มูลค่า":   fmt.Sprintf("฿%.2f", l.Amount),
					"ที่ตั้ง":  locationStr,
				}}
			}
		}
	}

	if len(idsByType[AssetTypeAccount]) > 0 {
		if res, _ := u.assetClient.GetBatchAccount(ctx, &assetPb.GetBatchIdsRequest{Ids: idsByType[AssetTypeAccount]}); res != nil {
			for _, a := range res.Account {
				infoMap[a.Id] = AssetInfo{Name: a.Name, Details: map[string]string{
					"ชื่อเรียก": a.Name,
					"ธนาคาร":    a.BankName,
					"เลขบัญชี":  utils.MaskBankAccount(a.BankAcc),
					"ยอดเงิน":   fmt.Sprintf("฿%.2f", a.Amount),
				}}
			}
		}
	}

	if len(idsByType[AssetTypeCash]) > 0 {
		if res, _ := u.assetClient.GetBatchCash(ctx, &assetPb.GetBatchIdsRequest{Ids: idsByType[AssetTypeCash]}); res != nil {
			for _, c := range res.Cash {
				infoMap[c.Id] = AssetInfo{Name: c.Name, Details: map[string]string{"ชื่อ": c.Name, "จำนวน": fmt.Sprintf("฿%.2f", c.Amount)}}
			}
		}
	}

	if len(idsByType[AssetTypeInsurance]) > 0 {
		if res, _ := u.assetClient.GetBatchInsurance(ctx, &assetPb.GetBatchIdsRequest{Ids: idsByType[AssetTypeInsurance]}); res != nil {
			for _, i := range res.Insurance {
				expDateStr := "ไม่ระบุ"
				if i.ExpDate != nil {
					expDateStr = i.ExpDate.AsTime().Format("02/01/2006")
				}

				infoMap[i.Id] = AssetInfo{Name: i.CompanyName, Details: map[string]string{
					"บริษัท":      i.CompanyName,
					"เลขกรมธรรม์": i.PolNum,
					"ทุนประกัน":   fmt.Sprintf("฿%.2f", i.CoverageAmount),
					"วันหมดอายุ":  expDateStr,
					"ประเภท":      i.Type.String(),
				}}
			}
		}
	}

	if len(idsByType[AssetTypeInvestment]) > 0 {
		if res, _ := u.assetClient.GetBatchInvestment(ctx, &assetPb.GetBatchIdsRequest{Ids: idsByType[AssetTypeInvestment]}); res != nil {
			for _, inv := range res.Invest {
				infoMap[inv.Id] = AssetInfo{Name: inv.Name, Details: map[string]string{
					"ชื่อกองทุน/หุ้น": inv.Name,
					"สัญลักษณ์":       inv.Symbol,
					"ประเภท":          inv.Type.String(),
				}}
			}
		}
	}

	if len(idsByType[AssetTypeLiability]) > 0 {
		if res, _ := u.assetClient.GetBatchLiability(ctx, &assetPb.GetBatchIdsRequest{Ids: idsByType[AssetTypeLiability]}); res != nil {
			for _, lia := range res.Liability {
				infoMap[lia.Id] = AssetInfo{Name: lia.Name, Details: map[string]string{
					"ชื่อหนี้สิน": lia.Name,
					"เจ้าหนี้":    lia.Creditor,
					"ยอดคงเหลือ":  fmt.Sprintf("฿%.2f", lia.Principal),
					"ประเภท":      lia.Type.String(),
				}}
			}
		}
	}

	for _, item := range items {
		data, found := infoMap[item.EntityID.String()]
		req := domain.SendEmailRequest{ToEmail: item.Email, AssetType: item.EntityType}

		if found {
			req.AssetName = data.Name
			req.ItemDetail = data.Details
		} else {
			req.AssetName = "รายการทรัพย์สิน"
		}

		if err := u.mailclient.SendShareInvitation(ctx, req); err != nil {
			fmt.Printf("❌ Email failed for %s: %v\n", item.Email, err)
		}
	}
}

func (u *ShareItemUsecase) ProcessScheduledEmails(ctx context.Context) error {
	pendingItems, err := u.itemRepo.GetPendingEmails(ctx)
	if err != nil {
		return err
	}

	if len(pendingItems) == 0 {
		return nil
	}

	var ids []uuid.UUID
	for _, item := range pendingItems {
		ids = append(ids, item.ID)
	}

	u.SendEmailInvitations(pendingItems)
	return u.itemRepo.MarkEmailsAsSent(ctx, ids)
}
