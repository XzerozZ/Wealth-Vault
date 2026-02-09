package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"wealth-vault/user-service/internal/domain"
	"wealth-vault/user-service/internal/event"
	repo "wealth-vault/user-service/internal/repository/interface"
	assetPb "wealth-vault/user-service/pkg/pb/proto/asset"
	pb "wealth-vault/user-service/pkg/pb/proto/user"
	"wealth-vault/user-service/pkg/utils"
	"wealth-vault/user-service/pkg/utils/mail"
	"wealth-vault/user-service/pkg/utils/mapper"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ShareItemUsecase struct {
	itemRepo    repo.ShareItemRepository
	groupRepo   repo.GroupRepository
	userRepo    repo.UserRepository
	msgRepo     repo.MsgRepository
	assetClient assetPb.AssetServiceClient
	mailclient  mail.NotificationClient
	publisher   *event.Publisher
}

func NewShareItemUsecase(
	r repo.ShareItemRepository,
	g repo.GroupRepository,
	u repo.UserRepository,
	m repo.MsgRepository,
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
		msgRepo:     m,
	}
}

func (u *ShareItemUsecase) ShareItem(ctx context.Context, req *pb.ShareItemRequest) (*pb.ShareItemResponse, error) {
	if len(req.ItemIds) == 0 {
		return nil, errors.New("no items to share")
	}
	if len(req.ItemIds) != len(req.ItemTypes) {
		return nil, errors.New("mismatch between item_ids and types length")
	}

	userID, _ := uuid.Parse(req.UserId)
	senderName := "Unknown"
	if userProfile, err := u.userRepo.GetUser(ctx, userID); err == nil {
		senderName = userProfile.Username
	}

	now := time.Now()
	var (
		groupItemsToCreate  []domain.GroupItem
		friendItemsToCreate []domain.FriendItem
		emailItemsToCreate  []domain.EmailItem
		emailsToSendNow     []domain.EmailItem
		groupLogsToCreate   []domain.GroupLog
		friendLogsToCreate  []domain.FriendLog
		groupActivities     []domain.GroupActivityEvent
		groupMsgsToCreate   []domain.GroupMessage
		privateMsgsToCreate []domain.PrivateMessage
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
		if err != nil || !res.Exists {
			return nil, fmt.Errorf("asset not found: %s", idStr)
		}

		assetDisplayName := fmt.Sprintf("%s shared item", entityType)
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

				logMeta := map[string]string{
					"action":    "share_item",
					"item_type": entityType,
					"item_id":   entityID.String(),
					"shared_at": shareTime.Format(time.RFC3339),
				}
				logMetaJSON, _ := json.Marshal(logMeta)

				groupLogsToCreate = append(groupLogsToCreate, domain.GroupLog{
					GroupID:   groupID,
					LogType:   "ACTIVITY",
					Messages:  fmt.Sprintf("%s ได้แชร์ %s รายการใหม่เข้ากลุ่ม", senderName, entityType),
					Metadata:  string(logMetaJSON),
					CreatedBy: userID,
				})

				groupActivities = append(groupActivities, domain.GroupActivityEvent{
					GroupID:      target.Id,
					ActivityType: "ITEM_SHARED",
					Payload:      fmt.Sprintf("%s แชร์ %s ใหม่", senderName, entityType),
					ActorID:      req.UserId,
					OccurredAt:   now.Unix(),
				})

				cardMeta := map[string]interface{}{
					"asset_id":       entityID.String(),
					"asset_type":     entityType,
					"snapshot_title": assetDisplayName,
					"action_url":     fmt.Sprintf("/asset/%s/%s", entityType, entityID),
				}

				cardMetaJSON, _ := json.Marshal(cardMeta)
				groupMsgsToCreate = append(groupMsgsToCreate, domain.GroupMessage{
					GroupID:   groupID,
					SenderID:  userID,
					MsgType:   "ASSET_CARD",
					Content:   "",
					Metadata:  string(cardMetaJSON),
					CreatedAt: now,
				})
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

				logMeta := map[string]string{
					"action":    "share_to_friend",
					"item_type": entityType,
					"item_id":   entityID.String(),
				}

				logMetaJSON, _ := json.Marshal(logMeta)
				friendLogsToCreate = append(friendLogsToCreate, domain.FriendLog{
					OwnerID:   userID,
					FriendID:  friendID,
					LogType:   "ACTIVITY",
					Messages:  fmt.Sprintf("คุณได้แชร์ %s รายการใหม่ให้กับเพื่อน", entityType),
					Metadata:  string(logMetaJSON),
					CreatedBy: userID,
				})

				cardMeta := map[string]interface{}{
					"asset_id":       entityID.String(),
					"asset_type":     entityType,
					"snapshot_title": assetDisplayName,
				}

				cardMetaJSON, _ := json.Marshal(cardMeta)
				privateMsgsToCreate = append(privateMsgsToCreate, domain.PrivateMessage{
					SenderID:   userID,
					ReceiverID: friendID,
					MsgType:    "ASSET_CARD",
					Content:    "",
					Metadata:   string(cardMetaJSON),
					CreatedAt:  now,
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

		go u.asyncSaveGroupLogs(groupLogsToCreate)
		go u.asyncBroadcastGroupActivities(groupActivities)
		go u.asyncSaveGroupMessages(groupMsgsToCreate)
	}

	if len(friendItemsToCreate) > 0 {
		if err := u.itemRepo.ShareItemtoFriend(ctx, friendItemsToCreate); err != nil {
			return nil, err
		}

		go u.asyncSaveFriendLogs(friendLogsToCreate)
		go u.asyncSavePrivateMessages(privateMsgsToCreate)
	}

	if len(emailItemsToCreate) > 0 {
		if err := u.itemRepo.ShareItemtoEmail(ctx, emailItemsToCreate); err != nil {
			return nil, err
		}
		if len(emailsToSendNow) > 0 {
			go u.SendEmailInvitations(emailsToSendNow)
		}
	}

	return &pb.ShareItemResponse{Finish: true}, nil
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
	groupID, _ := uuid.Parse(req.GroupId)
	uid, _ := uuid.Parse(req.UserId)

	items, err := u.itemRepo.GetSharedIteminGroup(ctx, groupID, uid)
	if err != nil {
		return nil, err
	}

	summaries := make([]domain.SharedItemSummary, len(items))
	for i, item := range items {
		summaries[i] = domain.SharedItemSummary{
			EntityID:   item.EntityID.String(),
			EntityType: item.EntityType,
		}
	}

	previewMap, err := u.fetchAssetPreviews(ctx, summaries)
	if err != nil {
		log.Printf("Failed to fetch asset previews: %v", err)
	}

	var responseItems []*pb.GroupItemDetail
	for _, item := range items {
		preview := previewMap[item.EntityID.String()]
		responseItems = append(responseItems, &pb.GroupItemDetail{
			GroupItemId: item.ID.String(),
			SharedBy:    item.OwnerID.String(),
			SharedAt:    timestamppb.New(item.CreatedAt),
			Type:        item.EntityType,
			AssetDetail: preview,
		})
	}

	return &pb.GetGroupItemsResponse{
		Items: responseItems,
	}, nil
}

func (u *ShareItemUsecase) GetSharedIteminFriend(ctx context.Context, req *pb.GetFriendItemRequest) (*pb.GetFriendItemsResponse, error) {
	uid, _ := uuid.Parse(req.UserId)
	fid, _ := uuid.Parse(req.FriendId)

	items, err := u.itemRepo.GetSharedIteminFriend(ctx, fid, uid)
	if err != nil {
		return nil, err
	}

	summaries := make([]domain.SharedItemSummary, len(items))
	for i, item := range items {
		summaries[i] = domain.SharedItemSummary{
			EntityID:   item.EntityID.String(),
			EntityType: item.EntityType,
		}
	}

	previewMap, err := u.fetchAssetPreviews(ctx, summaries)
	if err != nil {
		return nil, err
	}

	var responseItems []*pb.FriendItemDetail
	for _, item := range items {
		preview := previewMap[item.EntityID.String()]
		responseItems = append(responseItems, &pb.FriendItemDetail{
			FriendItemId: item.ID.String(),
			SharedBy:     item.OwnerID.String(),
			SharedAt:     timestamppb.New(item.CreatedAt),
			Type:         item.EntityType,
			AssetDetail:  preview,
		})
	}

	return &pb.GetFriendItemsResponse{
		Items: responseItems,
	}, nil
}

func (u *ShareItemUsecase) fetchAssetPreviews(ctx context.Context, items []domain.SharedItemSummary) (map[string]*pb.AssetPreview, error) {
	idsByType := make(map[string][]string)
	for _, item := range items {
		key := strings.ToLower(item.EntityType)
		idsByType[key] = append(idsByType[key], item.EntityID)
	}

	previewMap := make(map[string]*pb.AssetPreview)
	var wg sync.WaitGroup
	var mu sync.Mutex
	fetch := func(assetType string, fetcher func() error) {
		ids := idsByType[assetType]
		if len(ids) == 0 {
			return
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fetcher(); err != nil {
				log.Printf("⚠️ Failed to fetch %s assets: %v", assetType, err)
			}
		}()
	}

	fetch("building", func() error {
		res, err := u.assetClient.GetBatchBuilding(ctx, &assetPb.GetBatchIdsRequest{Ids: idsByType["building"]})
		if err != nil {
			return err
		}

		mu.Lock()
		defer mu.Unlock()
		if res != nil {
			for _, b := range res.Building {
				previewMap[b.Id] = mapper.MapBuildingToPreview(b)
			}
		}
		return nil
	})

	fetch("land", func() error {
		res, err := u.assetClient.GetBatchLand(ctx, &assetPb.GetBatchIdsRequest{Ids: idsByType["land"]})
		if err != nil {
			return err
		}

		mu.Lock()
		defer mu.Unlock()
		if res != nil {
			for _, l := range res.Land {
				previewMap[l.Id] = mapper.MapLandToPreview(l)
			}
		}
		return nil
	})

	fetch("account", func() error {
		res, err := u.assetClient.GetBatchAccount(ctx, &assetPb.GetBatchIdsRequest{Ids: idsByType["account"]})
		if err != nil {
			return err
		}

		mu.Lock()
		defer mu.Unlock()
		if res != nil {
			for _, a := range res.Account {
				mappedData := mapper.MapAccountToPreview(a)
				log.Printf("➡️ Mapping Account: ID='%s', ResultIsNil=%v", a.Id, mappedData == nil)

				previewMap[a.Id] = mappedData
			}
		}
		return nil
	})

	fetch("cash", func() error {
		res, err := u.assetClient.GetBatchCash(ctx, &assetPb.GetBatchIdsRequest{Ids: idsByType["cash"]})
		if err != nil {
			return err
		}

		mu.Lock()
		defer mu.Unlock()
		if res != nil {
			for _, c := range res.Cash {
				previewMap[c.Id] = mapper.MapCashToPreview(c)
			}
		}
		return nil
	})

	fetch("insurance", func() error {
		res, err := u.assetClient.GetBatchInsurance(ctx, &assetPb.GetBatchIdsRequest{Ids: idsByType["insurance"]})
		if err != nil {
			return err
		}

		mu.Lock()
		defer mu.Unlock()
		if res != nil {
			for _, i := range res.Insurance {
				previewMap[i.Id] = mapper.MapInsuranceToPreview(i)
			}
		}
		return nil
	})

	fetch("investment", func() error {
		res, err := u.assetClient.GetBatchInvestment(ctx, &assetPb.GetBatchIdsRequest{Ids: idsByType["investment"]})
		if err != nil {
			return err
		}

		mu.Lock()
		defer mu.Unlock()
		if res != nil {
			for _, inv := range res.Invest {
				previewMap[inv.Id] = mapper.MapInvestmentToPreview(inv)
			}
		}
		return nil
	})

	fetch("liability", func() error {
		res, err := u.assetClient.GetBatchLiability(ctx, &assetPb.GetBatchIdsRequest{Ids: idsByType["liability"]})
		if err != nil {
			return err
		}

		mu.Lock()
		defer mu.Unlock()
		if res != nil {
			for _, l := range res.Liability {
				previewMap[l.Id] = mapper.MapLiabilityToPreview(l)
			}
		}
		return nil
	})

	wg.Wait()

	return previewMap, nil
}

func (u *ShareItemUsecase) UnsharedIteminGroup(ctx context.Context, req *pb.UnshareItemRequest) (*pb.ShareItemResponse, error) {
	itemID, _ := uuid.Parse(req.ItemId)
	userID, _ := uuid.Parse(req.UserId)

	if err := u.itemRepo.DeleteIteminGroup(ctx, itemID, userID); err != nil {
		return nil, err
	}

	return &pb.ShareItemResponse{
		Finish: true,
	}, nil
}

func (u *ShareItemUsecase) UnsharedIteminFriend(ctx context.Context, req *pb.UnshareItemRequest) (*pb.ShareItemResponse, error) {
	itemID, _ := uuid.Parse(req.ItemId)
	userID, _ := uuid.Parse(req.UserId)

	if err := u.itemRepo.DeleteIteminFriend(ctx, itemID, userID); err != nil {
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
	groupID, _ := uuid.Parse(req.GroupId)
	senderID, _ := uuid.Parse(req.UserId)

	if len(req.TargetUserIds) == 0 {
		return nil, errors.New("no users specified")
	}

	senderName := "Unknown"
	if sender, err := u.userRepo.GetUser(ctx, senderID); err == nil {
		senderName = sender.Username
	}

	var newMembers []domain.GroupMember
	var targetUUIDs []uuid.UUID
	var addedNames []string

	for _, userIDStr := range req.TargetUserIds {
		uid, err := uuid.Parse(userIDStr)
		if err == nil {
			userName := "Someone"
			if user, err := u.userRepo.GetUser(ctx, uid); err == nil {
				userName = user.Username
			}
			addedNames = append(addedNames, userName)
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

	logMeta := map[string]interface{}{
		"action":      "add_member",
		"target_ids":  req.TargetUserIds,
		"added_count": len(newMembers),
	}
	logMetaJSON, _ := json.Marshal(logMeta)

	groupLog := domain.GroupLog{
		GroupID:   groupID,
		LogType:   "SYSTEM",
		Messages:  fmt.Sprintf("%s เพิ่มสมาชิกใหม่ %d คน", senderName, len(newMembers)),
		Metadata:  string(logMetaJSON),
		CreatedBy: senderID,
	}

	go func() {
		if err := u.groupRepo.CreateLog(context.Background(), &groupLog); err != nil {
			log.Printf("⚠️ Failed to create group log: %v", err)
		}
	}()

	msgContent := fmt.Sprintf("%s เพิ่ม %s เข้ากลุ่ม", senderName, strings.Join(addedNames, ", "))
	sysMsg := domain.GroupMessage{
		GroupID:   groupID,
		SenderID:  senderID,
		MsgType:   "SYSTEM_ALERT",
		Content:   msgContent,
		Metadata:  "{}",
		CreatedAt: time.Now(),
	}

	go func() {
		if err := u.msgRepo.CreateMessage(context.Background(), []domain.GroupMessage{sysMsg}); err != nil {
			log.Printf("⚠️ Failed to create system message: %v", err)
		}
	}()

	notifyTargetMap := make(map[string]bool)
	ownerUUIDs, err := u.itemRepo.GetItemOwnersInGroup(ctx, groupID)
	if err != nil {
		fmt.Printf("⚠️ Warning: Failed to fetch owners: %v\n", err)
	}
	for _, id := range ownerUUIDs {
		if id != senderID {
			notifyTargetMap[id.String()] = true
		}
	}

	var notifyTargetIDs []string
	for id := range notifyTargetMap {
		notifyTargetIDs = append(notifyTargetIDs, id)
	}

	evt := domain.GroupMemberAddedEvent{
		GroupID:       req.GroupId,
		SenderID:      req.UserId,
		AddedUserIDs:  req.TargetUserIds,
		TargetUserIDs: notifyTargetIDs,
		OccurredAt:    time.Now().Unix(),
	}
	go func() {
		if err := u.publisher.Publish("noti.group.member.added", evt); err != nil {
			fmt.Printf("⚠️ Failed to publish member added event: %v\n", err)
		}
	}()

	activityEvt := domain.GroupActivityEvent{
		GroupID:      req.GroupId,
		ActivityType: "MEMBER_ADDED",
		Payload:      fmt.Sprintf("%s เพิ่มสมาชิกใหม่", senderName),
		ActorID:      req.UserId,
		OccurredAt:   time.Now().Unix(),
	}
	go func() {
		if err := u.publisher.Publish("noti.group.activity", activityEvt); err != nil {
			log.Printf("⚠️ Failed to publish group activity: %v", err)
		}
	}()

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
			go u.itemRepo.BatchCreateViewers(context.Background(), newPermissions)
		}
	}

	return &pb.ActionResponse{
		Success: true,
	}, nil
}

func (u *ShareItemUsecase) GrantAccess(ctx context.Context, req *pb.GrantAccessRequest) (*pb.ActionResponse, error) {
	ownerID, _ := uuid.Parse(req.OwnerUserId)
	targetID, _ := uuid.Parse(req.TargetUserId)
	groupID, _ := uuid.Parse(req.GroupId)

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

	validItemIDs, err := u.itemRepo.GetOwnedItemIDs(ctx, req.GroupItemIds, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %v", err)
	}

	if len(validItemIDs) == 0 {
		return nil, errors.New("permission denied: you do not own any of the selected items")
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

	ownerName := "Unknown"
	targetName := "Member"
	if u1, err := u.userRepo.GetUser(ctx, ownerID); err == nil {
		ownerName = u1.Username
	}
	if u2, err := u.userRepo.GetUser(ctx, targetID); err == nil {
		targetName = u2.Username
	}

	logMeta := map[string]interface{}{
		"action":     "grant_access",
		"target_id":  req.TargetUserId,
		"item_ids":   req.GroupItemIds,
		"item_count": len(req.GroupItemIds),
	}
	logMetaJSON, _ := json.Marshal(logMeta)

	groupLog := domain.GroupLog{
		GroupID:   groupID,
		LogType:   "ACTIVITY",
		Messages:  fmt.Sprintf("%s ให้สิทธิ์ %s ดูรายการย้อนหลัง %d รายการ", ownerName, targetName, len(req.GroupItemIds)),
		Metadata:  string(logMetaJSON),
		CreatedBy: ownerID,
	}

	go func() {
		if err := u.groupRepo.CreateLog(context.Background(), &groupLog); err != nil {
			log.Printf("⚠️ Failed to create grant access log: %v", err)
		}
	}()

	evt := domain.AccessGrantedEvent{
		GroupID:      req.GroupId,
		GrantorID:    req.OwnerUserId,
		GrantorName:  ownerName,
		TargetUserID: req.TargetUserId,
		ItemCount:    len(req.GroupItemIds),
		OccurredAt:   time.Now().Unix(),
	}

	go func() {
		if err := u.publisher.Publish("noti.access.granted", evt); err != nil {
			log.Printf("⚠️ Failed to publish grant access event: %v", err)
		}
	}()

	activityEvt := domain.GroupActivityEvent{
		GroupID:      req.GroupId,
		ActivityType: "ACCESS_GRANTED",
		Payload:      fmt.Sprintf("คุณได้รับสิทธิ์ดูรายการเพิ่ม %d รายการ", len(req.GroupItemIds)),
		ActorID:      req.OwnerUserId,
		OccurredAt:   time.Now().Unix(),
	}

	go func() {
		if err := u.publisher.Publish("noti.group.activity", activityEvt); err != nil {
			log.Printf("⚠️ Failed to publish group activity: %v", err)
		}
	}()

	return &pb.ActionResponse{
		Success: true,
	}, nil
}

func (u *ShareItemUsecase) DeleteAllReferencesByEntityID(ctx context.Context, req *pb.DeleteByEntityRequest) (*pb.DeleteByEntityResponse, error) {
	id, _ := uuid.Parse(req.EntityId)

	if err := u.itemRepo.DeleteAllReferencesByEntityID(ctx, id); err != nil {
		return nil, err
	}

	return &pb.DeleteByEntityResponse{
		Success: true,
	}, nil
}

func (u *ShareItemUsecase) BatchShareAssets(ctx context.Context, req domain.BatchShareRequest) error {
	existingMap, err := u.itemRepo.GetExistingSharedMap(ctx, req.OwnerID, req.TargetID)
	if err != nil {
		return fmt.Errorf("failed to fetch existing shares: %v", err)
	}

	var newItems []domain.FriendItem
	now := time.Now()
	processIDs := func(ids []string, entityType string) {
		for _, idStr := range ids {
			if idStr == "" {
				continue
			}

			key := fmt.Sprintf("%s:%s", entityType, idStr)
			if !existingMap[key] {
				if uid, err := uuid.Parse(idStr); err == nil {
					newItems = append(newItems, domain.FriendItem{
						OwnerID:    req.OwnerID,
						FriendID:   req.TargetID,
						EntityType: entityType,
						EntityID:   uid,
						ShareAt:    now,
					})
				}
			}
		}
	}

	processIDs(req.AccountIDs, "account")
	processIDs(req.BuildingIDs, "building")
	processIDs(req.CashIDs, "cash")
	processIDs(req.InsuranceIDs, "insurance")
	processIDs(req.InvestmentIDs, "investment")
	processIDs(req.LandIDs, "land")
	processIDs(req.LiabilityIDs, "liability")

	if len(newItems) > 0 {
		if err := u.itemRepo.ShareItemtoFriend(ctx, newItems); err != nil {
			return fmt.Errorf("failed to batch insert items: %v", err)
		}

		go func() {
			bgCtx := context.Background()
			senderName := "Unknown"
			if userProfile, err := u.userRepo.GetUser(bgCtx, req.OwnerID); err == nil {
				senderName = userProfile.Username
			}

			evt := domain.ItemSharedEvent{
				AssetID:       "",
				SenderID:      req.OwnerID.String(),
				SenderName:    senderName,
				TargetUserIDs: []string{req.TargetID.String()},
				OccurredAt:    time.Now().Unix(),
			}

			u.publisher.Publish("noti.item.shared", evt)
		}()
	}

	return nil
}
