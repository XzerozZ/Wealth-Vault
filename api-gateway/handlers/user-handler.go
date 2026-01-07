package handlers

import (
	"context"
	"time"

	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/user"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	client pb.UserServiceClient
}

func NewUserHandler(c pb.UserServiceClient) *UserHandler {
	return &UserHandler{client: c}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var body domain.CreateUser
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := h.client.CreateUser(ctx, &pb.CreateUserRequest{
		Email:    body.Email,
		Username: body.Username,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}
