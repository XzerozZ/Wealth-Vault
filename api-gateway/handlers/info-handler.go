package handlers

import (
	"context"
	"time"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"
)

type InfoHandler struct {
	client pb.AssetServiceClient
}

func NewInfoHandler(c pb.AssetServiceClient) *InfoHandler {
	return &InfoHandler{
		client: c,
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
		assetsRes   *pb.GetAllAssetsResponse
		netWorthRes *pb.GetNetWorthResponse
	)

	g, gCtx := errgroup.WithContext(timeoutCtx)

	g.Go(func() error {
		var err error
		assetsRes, err = h.client.GetAllAssets(gCtx, &pb.GetAllAssetsRequest{
			UserId: userID,
		})
		return err
	})

	g.Go(func() error {
		var err error
		netWorthRes, err = h.client.GetNetWorth(gCtx, &pb.GetNetWorthRequest{
			UserId: userID,
		})
		return err
	})

	if err := g.Wait(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"assets": assetsRes.Assets,
		"net_worth": fiber.Map{
			"count":             netWorthRes.ItemCount,
			"total_assets":      netWorthRes.AssetsValue,
			"total_liabilities": netWorthRes.LiabilitiesValue,
			"value":             netWorthRes.NetWorth,
		},
	})
}
