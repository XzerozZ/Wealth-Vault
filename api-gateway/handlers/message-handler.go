package handlers

import (
	"context"
	"net/http"
	"time"
	"wealth-vault/api-gateway/internal/mapper"
	pb "wealth-vault/api-gateway/pkg/pb/proto/user"

	"github.com/gofiber/fiber/v2"
)

type MessageHandler struct {
	client pb.UserServiceClient
}

func NewMessageHandlerr(c pb.UserServiceClient) *MessageHandler {
	return &MessageHandler{
		client: c,
	}
}
func (h *MessageHandler) GetGroupMessages(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetGroupMessages(ctx, &pb.GetGroupMessagesRequest{
		GroupId: id,
		UserId:  userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	message := mapper.ToMessageResponseList(res.Messages)
	return c.JSON(fiber.Map{
		"messages": message,
	})
}

func (h *MessageHandler) GetPrivateMessages(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	req := &pb.GetPrivateMessagesRequest{
		UserId:   userID,
		FriendId: id,
	}

	res, err := h.client.GetPrivateMessages(c.Context(), req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	message := mapper.ToMessageResponseList(res.Messages)

	return c.JSON(fiber.Map{
		"messages": message,
	})
}
