package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"
	"wealth-vault/api-gateway/internal/domain"
	"wealth-vault/api-gateway/internal/mapper"
	assetpb "wealth-vault/api-gateway/pkg/pb/proto/asset"
	pb "wealth-vault/api-gateway/pkg/pb/proto/user"
	"wealth-vault/api-gateway/pkg/utils/helper"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"
)

type GroupItemHandler struct {
	client      pb.UserServiceClient
	assetClient assetpb.AssetServiceClient
}

func NewGroupItemHandler(c pb.UserServiceClient, ac assetpb.AssetServiceClient) *GroupItemHandler {
	return &GroupItemHandler{
		client:      c,
		assetClient: ac,
	}
}

func (h *GroupItemHandler) ShareItem(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.ShareItemRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if len(req.ItemIDs) != len(req.ItemTypes) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Mismatch between item_ids and item_types length",
		})
	}

	grpcReq := &pb.ShareItemRequest{
		UserId:    userID,
		ItemIds:   req.ItemIDs,
		ItemTypes: req.ItemTypes,
		Groups:    helper.MapShareTargets(req.Groups),
		Friends:   helper.MapShareTargets(req.Friends),
		Emails:    helper.MapShareTargets(req.Emails),
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.ShareItem(ctx, grpcReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res.Finish,
	})
}

func (h *GroupItemHandler) GetGroupItems(c *fiber.Ctx) error {
	groupID := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	res, err := h.client.GetSharedItem(c.Context(), &pb.GetGroupItemsRequest{
		GroupId: groupID,
		UserId:  userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	responseItems := mapper.MapGroupItemsToDomain(res.Items)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": responseItems,
	})
}

func (h *GroupItemHandler) GetFriendItems(c *fiber.Ctx) error {
	friendID := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	res, err := h.client.GetSharedIteminFriend(c.Context(), &pb.GetFriendItemRequest{
		FriendId: friendID,
		UserId:   userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	responseItems := mapper.MapFriendItemsToDomain(res.Items)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": responseItems,
	})
}

func (h *GroupItemHandler) UnsharedItem(c *fiber.Ctx) error {
	itemID := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	res, err := h.client.UnsharedItem(c.Context(), &pb.UnshareItemRequest{
		ItemId: itemID,
		UserId: userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res.Finish,
	})
}

func (h *GroupItemHandler) UnsharedIteminFriend(c *fiber.Ctx) error {
	itemID := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	res, err := h.client.UnsharedIteminFriend(c.Context(), &pb.UnshareItemRequest{
		ItemId: itemID,
		UserId: userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": res.Finish,
	})
}

func (h *GroupItemHandler) GetItemSharedTargets(c *fiber.Ctx) error {
	itemType := c.Params("type")
	itemID := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	res, err := h.client.GetItemSharedTargets(c.Context(), &pb.GetItemSharedTargetsRequest{
		UserId:   userID,
		ItemId:   itemID,
		ItemType: itemType,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	dtoResponse := mapper.ToSharedTargetsResponse(res)

	return c.JSON(dtoResponse)
}

func (h *GroupItemHandler) GetItemsForSelection(c *fiber.Ctx) error {
	targetID := c.Params("id")
	targetType := c.Params("type")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var (
		assetRes  *assetpb.GetAllAssetsResponse
		sharedRes *pb.GetSharedItemIDsResponse
	)

	g, ctx := errgroup.WithContext(c.Context())
	g.Go(func() error {
		var err error
		assetRes, err = h.assetClient.GetAllAssets(ctx, &assetpb.GetAllAssetsRequest{
			UserId: userID,
		})
		return err
	})

	g.Go(func() error {
		var err error
		sharedRes, err = h.client.GetSharedItemIDs(ctx, &pb.GetSharedItemIDsRequest{
			UserId:     userID,
			TargetId:   targetID,
			TargetType: targetType,
		})
		return err
	})

	if err := g.Wait(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data: " + err.Error()})
	}

	sharedMap := make(map[string]bool)
	for _, id := range sharedRes.ItemIds {
		sharedMap[strings.ToLower(id)] = true
	}

	var response []domain.AssetSelection
	for _, asset := range assetRes.Assets {
		isShared := false
		if _, exists := sharedMap[strings.ToLower(asset.Id)]; exists {
			isShared = true
		}

		response = append(response, domain.AssetSelection{
			ID:       asset.Id,
			Type:     asset.Type,
			Name:     asset.Name,
			Value:    asset.Value,
			Image:    asset.Image,
			IsShared: isShared,
		})
	}

	for _, asset := range assetRes.Liabilities {
		isShared := false
		if _, exists := sharedMap[strings.ToLower(asset.Id)]; exists {
			isShared = true
		}

		response = append(response, domain.AssetSelection{
			ID:       asset.Id,
			Type:     asset.Type,
			Name:     asset.Name,
			Value:    asset.Value,
			Image:    asset.Image,
			IsShared: isShared,
		})
	}

	return c.JSON(fiber.Map{
		"items": response,
	})
}
