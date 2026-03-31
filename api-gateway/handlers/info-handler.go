package handlers

import (
	"context"
	"sync"
	"time"
	assetPb "wealth-vault/api-gateway/pkg/pb/proto/asset"
	userPb "wealth-vault/api-gateway/pkg/pb/proto/user"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"
)

type InfoHandler struct {
	assetClient assetPb.AssetServiceClient
	userClient  userPb.UserServiceClient
}

func NewInfoHandler(ac assetPb.AssetServiceClient, uc userPb.UserServiceClient) *InfoHandler {
	return &InfoHandler{
		assetClient: ac,
		userClient:  uc,
	}
}

func (h *InfoHandler) Dashboard(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	timeoutCtx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	var (
		assetsRes     *assetPb.GetAllAssetsResponse
		netWorthRes   *assetPb.GetNetWorthResponse
		friendListRes *userPb.FriendListResponse
	)

	g, gCtx := errgroup.WithContext(timeoutCtx)
	g.Go(func() error {
		var err error
		assetsRes, err = h.assetClient.GetAllAssets(gCtx, &assetPb.GetAllAssetsRequest{
			UserId: userID,
		})
		return err
	})

	g.Go(func() error {
		var err error
		netWorthRes, err = h.assetClient.GetNetWorth(gCtx, &assetPb.GetNetWorthRequest{
			UserId: userID,
		})
		return err
	})

	g.Go(func() error {
		var err error
		friendListRes, err = h.userClient.GetFriendList(gCtx, &userPb.GetUserByIDRequest{
			Id: userID,
		})
		return err
	})

	if err := g.Wait(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	sharedItemMap := make(map[string]bool)
	if assetsRes != nil && assetsRes.Assets != nil {
		var checkGroup errgroup.Group
		var mu sync.Mutex

		for _, asset := range assetsRes.Assets {
			assetID := asset.Id
			assetType := asset.Type

			checkGroup.Go(func() error {
				targetRes, err := h.userClient.GetItemSharedTargets(c.UserContext(), &userPb.GetItemSharedTargetsRequest{
					UserId:   userID,
					ItemId:   assetID,
					ItemType: assetType,
				})

				if err == nil && targetRes != nil {
					if len(targetRes.Groups) > 0 || len(targetRes.Friends) > 0 || len(targetRes.Emails) > 0 {
						mu.Lock()
						sharedItemMap[assetID] = true
						mu.Unlock()
					}
				}

				return nil
			})
		}

		_ = checkGroup.Wait()
	}

	uniqueSharedItemCount := len(sharedItemMap)

	friendCount := 0
	if friendListRes != nil && friendListRes.Friends != nil {
		friendCount = len(friendListRes.Friends)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"assets":      assetsRes.Assets,
		"liabilities": assetsRes.Liabilities,
		"net_worth": fiber.Map{
			"count":             netWorthRes.ItemCount,
			"total_assets":      netWorthRes.AssetsValue,
			"total_liabilities": netWorthRes.LiabilitiesValue,
			"value":             netWorthRes.NetWorth,
		},
		"unique_shared_item_count": uniqueSharedItemCount,
		"friend_count":             friendCount,
	})
}
