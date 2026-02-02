package handlers

import (
	"context"
	"time"
	"wealth-vault/api-gateway/internal/domain"
	"wealth-vault/api-gateway/internal/mapper"
	pb "wealth-vault/api-gateway/pkg/pb/proto/user"
	"wealth-vault/api-gateway/pkg/utils/helper"

	"github.com/gofiber/fiber/v2"
)

type GroupItemHandler struct {
	client pb.UserServiceClient
}

func NewGroupItemHandler(c pb.UserServiceClient) *GroupItemHandler {
	return &GroupItemHandler{client: c}
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
