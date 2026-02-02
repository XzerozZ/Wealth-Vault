package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/event"
	repo "wealth-vault/user-service/internal/repository/interface"
	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils"
	"wealth-vault/user-service/pkg/utils/mail"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ShareItemUsecase struct {
	itemRepo    repo.ShareItemRepository
	groupRepo   repo.GroupRepository
	userRepo    repo.UserRepository
	assetClient assetPb.AssetServiceClient
	mailclient  mail.NotificationClient
	publisher   *event.Publisher
}

func NewShareItemUsecase(
	r repo.ShareItemRepository,
	g repo.GroupRepository,
	u repo.UserRepository,
	assetClient assetPb.AssetServiceClient,
	mail mail.NotificationClient,
	e *event.Publisher,
) ShareItemUsecase {
	return ShareItemUsecase{
		itemRepo:    r,
		groupRepo:   g,
		userRepo:    u,
		assetClient: assetClient,
		mailclient:  mail,
		publisher:   e,
	}
}

func (u *ShareItemUsecase) ShareItem(ctx context.Context, req *pb.ShareItemRequest) (*pb.ShareItemResponse, error) {
	if len(req.ItemIds) == 0 {
		return nil, errors.New("no items to share")
	}
	if len(req.ItemIds) != len(req.ItemTypes) {
		return nil, errors.New("mismatch between item_ids and types length")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	senderName := "Unknown"
	if userProfile, err := u.userRepo.GetUser(ctx, userID); err == nil {
		senderName = userProfile.Username
	}

	now := time.Now()

	var (
		groupItemsToCreate                  []domain.GroupItem
		friendItemsToCreate                 []domain.FriendItem
		emailItemsToCreate, emailsToSendNow []domain.EmailItem
	)

	for i, idStr := range req.ItemIds {
		assetNotifyMap := make(map[string]bool)
		entityType := req.ItemTypes[i]
		entityID, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid entity id at index %d: %v", i, err)
		}

		res, err := u.assetClient.CheckAssetExists(ctx, &assetPb.CheckAssetRequest{
			Id:     idStr,
			UserId: req.UserId,
			Type:   entityType,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to verify asset %s: %v", idStr, err)
		}
		if !res.Exists {
			return nil, fmt.Errorf("asset not found: %s", idStr)
		}

		for _, target := range req.Groups {
			groupID, _ := uuid.Parse(target.Id)

			shareTime := now
			if target.ShareAt != nil {
				shareTime = target.ShareAt.AsTime()
			}

			exist, _ := u.itemRepo.IsItemSharedtoGroup(ctx, groupID, entityID, entityType)
			if !exist {
				newItem := domain.GroupItem{
					GroupID:    groupID,
					EntityType: entityType,
					EntityID:   entityID,
					OwnerID:    userID,
					ShareAt:    shareTime,
				}

				members, _, err := u.groupRepo.GetMember(ctx, groupID)
				if err == nil {
					var viewers []domain.GroupItemViewer
					for _, member := range members {
						viewers = append(viewers, domain.GroupItemViewer{
							GroupItemID: newItem.ID,
							ViewerID:    member.ID,
						})

						if member.ID != userID {
							assetNotifyMap[member.ID.String()] = true
						}
					}

					newItem.Viewers = viewers
				}

				groupItemsToCreate = append(groupItemsToCreate, newItem)
			}
		}

		for _, target := range req.Friends {
			friendID, _ := uuid.Parse(target.Id)

			shareTime := now
			if target.ShareAt != nil {
				shareTime = target.ShareAt.AsTime()
			}

			exist, _ := u.itemRepo.IsItemSharedtoFriend(ctx, friendID, entityID, entityType)
			if !exist {
				friendItemsToCreate = append(friendItemsToCreate, domain.FriendItem{
					OwnerID:    userID,
					FriendID:   friendID,
					EntityType: entityType,
					EntityID:   entityID,
					ShareAt:    shareTime,
				})
			}

			assetNotifyMap[target.Id] = true
		}

		for _, target := range req.Emails {
			shareTime := now
			if target.ShareAt != nil {
				shareTime = target.ShareAt.AsTime()
			}

			shouldSendNow := shareTime.IsZero() || shareTime.Before(time.Now().Add(1*time.Minute))

			emailItem := domain.EmailItem{
				OwnerID:    userID,
				Email:      target.Id,
				EntityType: entityType,
				EntityID:   entityID,
				ShareAt:    shareTime,
				IsSent:     shouldSendNow,
			}

			emailItemsToCreate = append(emailItemsToCreate, emailItem)

			if shouldSendNow {
				emailsToSendNow = append(emailsToSendNow, emailItem)
			}
		}

		if len(assetNotifyMap) > 0 {
			var targetIDs []string
			for uid := range assetNotifyMap {
				targetIDs = append(targetIDs, uid)
			}

			evt := domain.ItemSharedEvent{
				SenderID:      req.UserId,
				SenderName:    senderName,
				AssetID:       idStr,
				TargetUserIDs: targetIDs,
				OccurredAt:    time.Now().Unix(),
			}

			go u.publisher.Publish("noti.item.shared", evt)
		}
	}

	if len(groupItemsToCreate) > 0 {
		if err := u.itemRepo.ShareItemtoGroup(ctx, groupItemsToCreate); err != nil {
			return nil, err
		}
	}

	if len(friendItemsToCreate) > 0 {
		if err := u.itemRepo.ShareItemtoFriend(ctx, friendItemsToCreate); err != nil {
			return nil, err
		}
	}

	if len(emailItemsToCreate) > 0 {
		if err := u.itemRepo.ShareItemtoEmail(ctx, emailItemsToCreate); err != nil {
			return nil, err
		}

		if len(emailsToSendNow) > 0 {
			go u.SendEmailInvitations(emailsToSendNow)
		}
	}

	return &pb.ShareItemResponse{
		Finish: true,
	}, nil
}

func (u *ShareItemUsecase) SendEmailInvitations(items []domain.EmailItem) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var buildingIDs, landIDs, accountIDs, cashIDs, insuranceIDs, investmentIDs, liabilityIDs []string

	for _, item := range items {
		id := item.EntityID.String()
		switch item.EntityType {
		case "building":
			buildingIDs = append(buildingIDs, id)
		case "land":
			landIDs = append(landIDs, id)
		case "account":
			accountIDs = append(accountIDs, id)
		case "cash":
			cashIDs = append(cashIDs, id)
		case "insurance":
			insuranceIDs = append(insuranceIDs, id)
		case "investment":
			investmentIDs = append(investmentIDs, id)
		case "liability":
			liabilityIDs = append(liabilityIDs, id)
		}
	}

	type AssetInfo struct {
		Name    string
		Image   string
		Details map[string]string
	}
	infoMap := make(map[string]AssetInfo)

	if len(buildingIDs) > 0 {
		res, _ := u.assetClient.GetBatchBuilding(ctx, &assetPb.GetBatchIdsRequest{Ids: buildingIDs})
		if res != nil {
			for _, b := range res.Building {
				infoMap[b.Id] = AssetInfo{
					Name: b.Name,
					Details: map[string]string{
						"ชื่อ":    b.Name,
						"มูลค่า":  fmt.Sprintf("฿%.2f", b.Amount),
						"ที่ตั้ง": fmt.Sprintf("%s, %s", b.Location.District, b.Location.Province), //
						"ประเภท":  b.Type.String(),
					},
				}
			}
		}
	}

	if len(landIDs) > 0 {
		res, _ := u.assetClient.GetBatchLand(ctx, &assetPb.GetBatchIdsRequest{Ids: landIDs})
		if res != nil {
			for _, l := range res.Land {
				infoMap[l.Id] = AssetInfo{
					Name: l.Name,
					Details: map[string]string{
						"ชื่อ":     l.Name,
						"เลขโฉนด":  l.DeedNum,
						"เนื้อที่": fmt.Sprintf("%.2f", l.Area),
						"มูลค่า":   fmt.Sprintf("฿%.2f", l.Amount),
						"ที่ตั้ง":  fmt.Sprintf("%s, %s", l.Location.District, l.Location.Province),
					},
				}
			}
		}
	}

	if len(accountIDs) > 0 {
		res, _ := u.assetClient.GetBatchAccount(ctx, &assetPb.GetBatchIdsRequest{Ids: accountIDs})
		if res != nil {
			for _, a := range res.Account {
				infoMap[a.Id] = AssetInfo{
					Name: a.Name,
					Details: map[string]string{
						"ชื่อเรียก": a.Name,
						"ธนาคาร":    a.BankName,
						"เลขบัญชี":  utils.MaskBankAccount(a.BankAcc),
						"ยอดเงิน":   fmt.Sprintf("฿%.2f", a.Amount),
					},
				}
			}
		}
	}

	if len(cashIDs) > 0 {
		res, _ := u.assetClient.GetBatchCash(ctx, &assetPb.GetBatchIdsRequest{Ids: cashIDs})
		if res != nil {
			for _, c := range res.Cash {
				infoMap[c.Id] = AssetInfo{
					Name: c.Name,
					Details: map[string]string{
						"ชื่อ":  c.Name,
						"จำนวน": fmt.Sprintf("฿%.2f", c.Amount),
					},
				}
			}
		}
	}

	if len(insuranceIDs) > 0 {
		res, _ := u.assetClient.GetBatchInsurance(ctx, &assetPb.GetBatchIdsRequest{Ids: insuranceIDs})
		if res != nil {
			for _, i := range res.Insurance {
				infoMap[i.Id] = AssetInfo{
					Name: i.CompanyName,
					Details: map[string]string{
						"บริษัท":      i.CompanyName,
						"เลขกรมธรรม์": i.PolNum,
						"ทุนประกัน":   fmt.Sprintf("฿%.2f", i.CoverageAmount),
						"วันหมดอายุ":  i.ExpDate.AsTime().Format("02/01/2006"),
						"ประเภท":      i.Type.String(),
					},
				}
			}
		}
	}

	if len(investmentIDs) > 0 {
		res, _ := u.assetClient.GetBatchInvestment(ctx, &assetPb.GetBatchIdsRequest{Ids: investmentIDs})
		if res != nil {
			for _, inv := range res.Invest {
				infoMap[inv.Id] = AssetInfo{
					Name: inv.Name,
					Details: map[string]string{
						"ชื่อกองทุน/หุ้น": inv.Name,
						"สัญลักษณ์":       inv.Symbol,
						"ประเภท":          inv.Type.String(),
					},
				}
			}
		}
	}

	if len(liabilityIDs) > 0 {
		res, _ := u.assetClient.GetBatchLiability(ctx, &assetPb.GetBatchIdsRequest{Ids: liabilityIDs})
		if res != nil {
			for _, lia := range res.Liability {
				infoMap[lia.Id] = AssetInfo{
					Name: lia.Name,
					Details: map[string]string{
						"ชื่อหนี้สิน": lia.Name,
						"เจ้าหนี้":    lia.Creditor,
						"ยอดคงเหลือ":  fmt.Sprintf("฿%.2f", lia.Principal),
						"ประเภท":      lia.Type.String(),
					},
				}
			}
		}
	}

	for _, item := range items {
		data, found := infoMap[item.EntityID.String()]

		req := domain.SendEmailRequest{
			ToEmail:   item.Email,
			AssetType: item.EntityType,
		}

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

func (u *ShareItemUsecase) GetSharedIteminGroup(ctx context.Context, req *pb.GetGroupItemsRequest) (*pb.GetGroupItemsResponse, error) {
	groupID, err := uuid.Parse(req.GroupId)
	if err != nil {
		return nil, errors.New("invalid group id")
	}

	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	items, err := u.itemRepo.GetSharedIteminGroup(ctx, groupID, uid)
	if err != nil {
		return nil, err
	}

	var buildingIDs, landIDs, accountIDs, cashIDs, insuranceIDs, investmentIDs, liabilityIDs []string

	itemMap := make(map[string]*domain.GroupItem)

	for _, item := range items {
		idStr := item.EntityID.String()
		itemMap[idStr] = &item

		switch item.EntityType {
		case "building":
			buildingIDs = append(buildingIDs, idStr)
		case "land":
			landIDs = append(landIDs, idStr)
		case "account":
			accountIDs = append(accountIDs, idStr)
		case "cash":
			cashIDs = append(cashIDs, idStr)
		case "insurance":
			insuranceIDs = append(insuranceIDs, idStr)
		case "investment":
			investmentIDs = append(investmentIDs, idStr)
		case "liability":
			liabilityIDs = append(liabilityIDs, idStr)
		}
	}

	var buildings []*assetPb.Building
	if len(buildingIDs) > 0 {
		res, err := u.assetClient.GetBatchBuilding(ctx, &assetPb.GetBatchIdsRequest{
			Ids: buildingIDs,
		})
		if err == nil && res != nil {
			buildings = res.Building
		}
	}

	var lands []*assetPb.Land
	if len(landIDs) > 0 {
		res, err := u.assetClient.GetBatchLand(ctx, &assetPb.GetBatchIdsRequest{
			Ids: landIDs,
		})
		if err == nil && res != nil {
			lands = res.Land
		}
	}

	var accounts []*assetPb.Account
	if len(accountIDs) > 0 {
		res, err := u.assetClient.GetBatchAccount(ctx, &assetPb.GetBatchIdsRequest{
			Ids: accountIDs,
		})
		if err == nil && res != nil {
			accounts = res.Account
		}
	}

	var cashes []*assetPb.Cash
	if len(cashIDs) > 0 {
		res, err := u.assetClient.GetBatchCash(ctx, &assetPb.GetBatchIdsRequest{
			Ids: cashIDs,
		})
		if err == nil && res != nil {
			cashes = res.Cash
		}
	}

	var insurances []*assetPb.Insurance
	if len(insuranceIDs) > 0 {
		res, err := u.assetClient.GetBatchInsurance(ctx, &assetPb.GetBatchIdsRequest{
			Ids: insuranceIDs,
		})
		if err == nil && res != nil {
			insurances = res.Insurance
		}
	}

	var investments []*assetPb.Investment
	if len(investmentIDs) > 0 {
		res, err := u.assetClient.GetBatchInvestment(ctx, &assetPb.GetBatchIdsRequest{
			Ids: investmentIDs,
		})
		if err == nil && res != nil {
			investments = res.Invest
		}
	}

	var liabilities []*assetPb.Liability
	if len(liabilityIDs) > 0 {
		res, err := u.assetClient.GetBatchLiability(ctx, &assetPb.GetBatchIdsRequest{
			Ids: liabilityIDs,
		})
		if err == nil && res != nil {
			liabilities = res.Liability
		}
	}

	var responseItems []*pb.GroupItemDetail
	createBaseDetail := func(assetID string) *pb.GroupItemDetail {
		if origin, ok := itemMap[assetID]; ok {
			return &pb.GroupItemDetail{
				GroupItemId: origin.ID.String(),
				SharedBy:    origin.OwnerID.String(),
				SharedAt:    timestamppb.New(origin.CreatedAt),
			}
		}
		return nil
	}

	for _, b := range buildings {
		if detail := createBaseDetail(b.Id); detail != nil {
			detail.AssetDetail = &pb.AssetPreview{
				Asset: &pb.AssetPreview_Building{
					Building: &pb.BuildingPreview{
						Id:           b.Id,
						Name:         b.Name,
						Amount:       b.Amount,
						LocationText: fmt.Sprintf("%s, %s", b.Location.District, b.Location.Province),
						TypeName:     b.Type.String(),
					},
				},
			}

			responseItems = append(responseItems, detail)
		}
	}

	for _, l := range lands {
		if detail := createBaseDetail(l.Id); detail != nil {
			detail.AssetDetail = &pb.AssetPreview{
				Asset: &pb.AssetPreview_Land{
					Land: &pb.LandPreview{
						Id:           l.Id,
						Name:         l.Name,
						DeedNum:      l.DeedNum,
						Area:         l.Area,
						Amount:       l.Amount,
						LocationText: fmt.Sprintf("%s, %s", l.Location.District, l.Location.Province),
					},
				},
			}

			responseItems = append(responseItems, detail)
		}
	}

	for _, a := range accounts {
		if detail := createBaseDetail(a.Id); detail != nil {
			detail.AssetDetail = &pb.AssetPreview{
				Asset: &pb.AssetPreview_Account{
					Account: &pb.AccountPreview{
						Id:            a.Id,
						Name:          a.Name,
						BankName:      a.BankName,
						Amount:        a.Amount,
						AccountNumber: utils.MaskBankAccount(a.BankAcc),
					},
				},
			}

			responseItems = append(responseItems, detail)
		}
	}

	for _, c := range cashes {
		if detail := createBaseDetail(c.Id); detail != nil {
			detail.AssetDetail = &pb.AssetPreview{
				Asset: &pb.AssetPreview_Cash{
					Cash: &pb.CashPreview{
						Id:     c.Id,
						Name:   c.Name,
						Amount: c.Amount,
					},
				},
			}

			responseItems = append(responseItems, detail)
		}
	}

	for _, i := range insurances {
		if detail := createBaseDetail(i.Id); detail != nil {
			detail.AssetDetail = &pb.AssetPreview{
				Asset: &pb.AssetPreview_Insurance{
					Insurance: &pb.InsurancePreview{
						Id:             i.Id,
						TypeName:       i.Type.String(),
						CompanyName:    i.CompanyName,
						PolNum:         i.PolNum,
						CoverageAmount: i.CoverageAmount,
						ExpDateText:    i.ExpDate.AsTime().String(),
					},
				},
			}

			responseItems = append(responseItems, detail)
		}
	}

	for _, inv := range investments {
		if detail := createBaseDetail(inv.Id); detail != nil {
			detail.AssetDetail = &pb.AssetPreview{
				Asset: &pb.AssetPreview_Investment{
					Investment: &pb.InvestmentPreview{
						Id:       inv.Id,
						Name:     inv.Name,
						TypeName: inv.Type.String(),
						Symbol:   inv.Symbol,
					},
				},
			}

			responseItems = append(responseItems, detail)
		}
	}

	for _, lia := range liabilities {
		if detail := createBaseDetail(lia.Id); detail != nil {
			detail.AssetDetail = &pb.AssetPreview{
				Asset: &pb.AssetPreview_Liability{
					Liability: &pb.LiabilityPreview{
						Id:        lia.Id,
						Name:      lia.Name,
						TypeName:  lia.Type.String(),
						Creditor:  lia.Creditor,
						Principal: lia.Principal,
					},
				},
			}

			responseItems = append(responseItems, detail)
		}
	}

	return &pb.GetGroupItemsResponse{
		Items: responseItems,
	}, nil
}

func (u *ShareItemUsecase) GetSharedIteminFriend(ctx context.Context, req *pb.GetFriendItemRequest) (*pb.GetFriendItemsResponse, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	fid, err := uuid.Parse(req.FriendId)
	if err != nil {
		return nil, errors.New("invalid friend id")
	}

	items, err := u.itemRepo.GetSharedIteminFriend(ctx, fid, uid)
	if err != nil {
		return nil, err
	}

	var buildingIDs, landIDs, accountIDs, cashIDs, insuranceIDs, investmentIDs, liabilityIDs []string

	itemMap := make(map[string]*domain.FriendItem)

	for _, item := range items {
		idStr := item.EntityID.String()
		itemMap[idStr] = &item

		switch item.EntityType {
		case "building":
			buildingIDs = append(buildingIDs, idStr)
		case "land":
			landIDs = append(landIDs, idStr)
		case "account":
			accountIDs = append(accountIDs, idStr)
		case "cash":
			cashIDs = append(cashIDs, idStr)
		case "insurance":
			insuranceIDs = append(insuranceIDs, idStr)
		case "investment":
			investmentIDs = append(investmentIDs, idStr)
		case "liability":
			liabilityIDs = append(liabilityIDs, idStr)
		}
	}

	var buildings []*assetPb.Building
	if len(buildingIDs) > 0 {
		res, err := u.assetClient.GetBatchBuilding(ctx, &assetPb.GetBatchIdsRequest{
			Ids: buildingIDs,
		})
		if err == nil && res != nil {
			buildings = res.Building
		}
	}

	var lands []*assetPb.Land
	if len(landIDs) > 0 {
		res, err := u.assetClient.GetBatchLand(ctx, &assetPb.GetBatchIdsRequest{
			Ids: landIDs,
		})
		if err == nil && res != nil {
			lands = res.Land
		}
	}

	var accounts []*assetPb.Account
	if len(accountIDs) > 0 {
		res, err := u.assetClient.GetBatchAccount(ctx, &assetPb.GetBatchIdsRequest{
			Ids: accountIDs,
		})
		if err == nil && res != nil {
			accounts = res.Account
		}
	}

	var cashes []*assetPb.Cash
	if len(cashIDs) > 0 {
		res, err := u.assetClient.GetBatchCash(ctx, &assetPb.GetBatchIdsRequest{
			Ids: cashIDs,
		})
		if err == nil && res != nil {
			cashes = res.Cash
		}
	}

	var insurances []*assetPb.Insurance
	if len(insuranceIDs) > 0 {
		res, err := u.assetClient.GetBatchInsurance(ctx, &assetPb.GetBatchIdsRequest{
			Ids: insuranceIDs,
		})
		if err == nil && res != nil {
			insurances = res.Insurance
		}
	}

	var investments []*assetPb.Investment
	if len(investmentIDs) > 0 {
		res, err := u.assetClient.GetBatchInvestment(ctx, &assetPb.GetBatchIdsRequest{
			Ids: investmentIDs,
		})
		if err == nil && res != nil {
			investments = res.Invest
		}
	}

	var liabilities []*assetPb.Liability
	if len(liabilityIDs) > 0 {
		res, err := u.assetClient.GetBatchLiability(ctx, &assetPb.GetBatchIdsRequest{
			Ids: liabilityIDs,
		})
		if err == nil && res != nil {
			liabilities = res.Liability
		}
	}

	var responseItems []*pb.FriendItemDetail
	createBaseDetail := func(assetID string) *pb.FriendItemDetail {
		if origin, ok := itemMap[assetID]; ok {
			return &pb.FriendItemDetail{
				FriendItemId: origin.ID.String(),
				SharedBy:     origin.OwnerID.String(),
				SharedAt:     timestamppb.New(origin.CreatedAt),
			}
		}
		return nil
	}

	for _, b := range buildings {
		if detail := createBaseDetail(b.Id); detail != nil {
			detail.AssetDetail = &pb.AssetPreview{
				Asset: &pb.AssetPreview_Building{
					Building: &pb.BuildingPreview{
						Id:           b.Id,
						Name:         b.Name,
						Amount:       b.Amount,
						LocationText: fmt.Sprintf("%s, %s", b.Location.District, b.Location.Province),
						TypeName:     b.Type.String(),
					},
				},
			}

			responseItems = append(responseItems, detail)
		}
	}

	for _, l := range lands {
		if detail := createBaseDetail(l.Id); detail != nil {
			detail.AssetDetail = &pb.AssetPreview{
				Asset: &pb.AssetPreview_Land{
					Land: &pb.LandPreview{
						Id:           l.Id,
						Name:         l.Name,
						DeedNum:      l.DeedNum,
						Area:         l.Area,
						Amount:       l.Amount,
						LocationText: fmt.Sprintf("%s, %s", l.Location.District, l.Location.Province),
					},
				},
			}

			responseItems = append(responseItems, detail)
		}
	}

	for _, a := range accounts {
		if detail := createBaseDetail(a.Id); detail != nil {
			detail.AssetDetail = &pb.AssetPreview{
				Asset: &pb.AssetPreview_Account{
					Account: &pb.AccountPreview{
						Id:            a.Id,
						Name:          a.Name,
						BankName:      a.BankName,
						Amount:        a.Amount,
						AccountNumber: utils.MaskBankAccount(a.BankAcc),
					},
				},
			}

			responseItems = append(responseItems, detail)
		}
	}

	for _, c := range cashes {
		if detail := createBaseDetail(c.Id); detail != nil {
			detail.AssetDetail = &pb.AssetPreview{
				Asset: &pb.AssetPreview_Cash{
					Cash: &pb.CashPreview{
						Id:     c.Id,
						Name:   c.Name,
						Amount: c.Amount,
					},
				},
			}

			responseItems = append(responseItems, detail)
		}
	}

	for _, i := range insurances {
		if detail := createBaseDetail(i.Id); detail != nil {
			detail.AssetDetail = &pb.AssetPreview{
				Asset: &pb.AssetPreview_Insurance{
					Insurance: &pb.InsurancePreview{
						Id:             i.Id,
						TypeName:       i.Type.String(),
						CompanyName:    i.CompanyName,
						PolNum:         i.PolNum,
						CoverageAmount: i.CoverageAmount,
						ExpDateText:    i.ExpDate.AsTime().String(),
					},
				},
			}

			responseItems = append(responseItems, detail)
		}
	}

	for _, inv := range investments {
		if detail := createBaseDetail(inv.Id); detail != nil {
			detail.AssetDetail = &pb.AssetPreview{
				Asset: &pb.AssetPreview_Investment{
					Investment: &pb.InvestmentPreview{
						Id:       inv.Id,
						Name:     inv.Name,
						TypeName: inv.Type.String(),
						Symbol:   inv.Symbol,
					},
				},
			}

			responseItems = append(responseItems, detail)
		}
	}

	for _, lia := range liabilities {
		if detail := createBaseDetail(lia.Id); detail != nil {
			detail.AssetDetail = &pb.AssetPreview{
				Asset: &pb.AssetPreview_Liability{
					Liability: &pb.LiabilityPreview{
						Id:        lia.Id,
						Name:      lia.Name,
						TypeName:  lia.Type.String(),
						Creditor:  lia.Creditor,
						Principal: lia.Principal,
					},
				},
			}

			responseItems = append(responseItems, detail)
		}
	}

	return &pb.GetFriendItemsResponse{
		Items: responseItems,
	}, nil
}

func (u *ShareItemUsecase) UnsharedIteminGroup(ctx context.Context, req *pb.UnshareItemRequest) (*pb.ShareItemResponse, error) {
	itemID, err := uuid.Parse(req.ItemId)
	if err != nil {
		return nil, errors.New("invalid item id")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	if err = u.itemRepo.DeleteIteminGroup(ctx, itemID, userID); err != nil {
		return nil, err
	}

	return &pb.ShareItemResponse{
		Finish: true,
	}, nil
}

func (u *ShareItemUsecase) UnsharedIteminFriend(ctx context.Context, req *pb.UnshareItemRequest) (*pb.ShareItemResponse, error) {
	itemID, err := uuid.Parse(req.ItemId)
	if err != nil {
		return nil, errors.New("invalid item id")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	if err = u.itemRepo.DeleteIteminFriend(ctx, itemID, userID); err != nil {
		return nil, err
	}

	return &pb.ShareItemResponse{
		Finish: true,
	}, nil
}

func (u *ShareItemUsecase) ProcessScheduledEmails(ctx context.Context) error {
	pendingItems, err := u.itemRepo.GetPendingEmails(ctx)
	if err != nil {
		return err
	}

	if len(pendingItems) == 0 {
		fmt.Println("No pending emails to send.")
		return nil
	}

	var ids []uuid.UUID
	for _, item := range pendingItems {
		ids = append(ids, item.ID)
	}

	u.SendEmailInvitations(pendingItems)

	if err := u.itemRepo.MarkEmailsAsSent(ctx, ids); err != nil {
		return err
	}

	return nil
}

func (u *ShareItemUsecase) AddMemberToGroup(ctx context.Context, req *pb.AddMemberRequest) (*pb.ActionResponse, error) {
	groupID, err := uuid.Parse(req.GroupId)
	if err != nil {
		return nil, errors.New("invalid group id")
	}

	if len(req.TargetUserIds) == 0 {
		return nil, errors.New("no users specified")
	}

	var newMembers []domain.GroupMember
	var targetUUIDs []uuid.UUID
	for _, userIDStr := range req.TargetUserIds {
		uid, err := uuid.Parse(userIDStr)
		if err == nil {
			targetUUIDs = append(targetUUIDs, uid)
			newMembers = append(newMembers, domain.GroupMember{
				GroupID:  groupID,
				UserID:   uid,
				Role:     "member",
				JoinedAt: time.Now(),
			})
		}
	}

	if len(newMembers) > 0 {
		if err := u.itemRepo.AddMember(ctx, newMembers); err != nil {
			return nil, fmt.Errorf("failed to add members: %v", err)
		}
	}

	ownerUUIDs, err := u.itemRepo.GetItemOwnersInGroup(ctx, groupID)
	if err != nil {
		fmt.Printf("⚠️ Warning: Failed to fetch owners: %v\n", err)
	}

	var ownerIDStrings []string
	for _, id := range ownerUUIDs {
		ownerIDStrings = append(ownerIDStrings, id.String())
	}

	evt := domain.GroupMemberAddedEvent{
		GroupID:       req.GroupId,
		SenderID:      req.UserId,
		AddedUserIDs:  req.TargetUserIds,
		TargetUserIDs: ownerIDStrings,
		OccurredAt:    time.Now().Unix(),
	}

	if err := u.publisher.Publish(event.GroupMemberAdded, evt); err != nil {
		fmt.Printf("⚠️ Failed to publish event: %v\n", err)
	}

	futureItemIDs, err := u.itemRepo.GetFutureItemsInGroup(ctx, groupID)
	if err != nil {
		fmt.Printf("⚠️ Warning: Failed to get future items for auto-grant: %v\n", err)
	}

	if len(futureItemIDs) > 0 && len(targetUUIDs) > 0 {
		var newPermissions []domain.GroupItemViewer
		for _, userID := range targetUUIDs {
			for _, itemID := range futureItemIDs {
				newPermissions = append(newPermissions, domain.GroupItemViewer{
					GroupItemID: itemID,
					ViewerID:    userID,
					GrantedAt:   time.Now(),
				})
			}
		}

		if len(newPermissions) > 0 {
			_ = u.itemRepo.BatchCreateViewers(ctx, newPermissions)
		}
	}

	return &pb.ActionResponse{
		Success: true,
	}, nil
}

func (u *ShareItemUsecase) GrantAccess(ctx context.Context, req *pb.GrantAccessRequest) (*pb.ActionResponse, error) {
	ownerID, err := uuid.Parse(req.OwnerUserId)
	if err != nil {
		return nil, errors.New("invalid owner id")
	}
	targetID, err := uuid.Parse(req.TargetUserId)
	if err != nil {
		return nil, errors.New("invalid target user id")
	}
	groupID, err := uuid.Parse(req.GroupId)
	if err != nil {
		return nil, errors.New("invalid group id")
	}

	if len(req.GroupItemIds) == 0 {
		return nil, errors.New("no items selected")
	}

	isMember, err := u.itemRepo.IsGroupMember(ctx, groupID, targetID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("target user is not a member of this group")
	}

	count, err := u.itemRepo.CountItemsByOwner(ctx, req.GroupItemIds, ownerID)
	if err != nil {
		return nil, err
	}
	if count != int64(len(req.GroupItemIds)) {
		return nil, errors.New("permission denied: you do not own all selected items or some items do not exist")
	}

	var viewers []domain.GroupItemViewer
	for _, itemIDStr := range req.GroupItemIds {
		itemID, _ := uuid.Parse(itemIDStr)
		viewers = append(viewers, domain.GroupItemViewer{
			GroupItemID: itemID,
			ViewerID:    targetID,
			GrantedAt:   time.Now(),
		})
	}

	if err := u.itemRepo.BatchCreateViewers(ctx, viewers); err != nil {
		return nil, fmt.Errorf("failed to grant access: %v", err)
	}

	return &pb.ActionResponse{
		Success: true,
	}, nil
}
