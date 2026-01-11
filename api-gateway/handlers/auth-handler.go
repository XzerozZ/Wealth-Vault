package handlers

import (
	"context"
	"time"

	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/auth"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	client pb.AuthServiceClient
}

func NewAuthHandler(c pb.AuthServiceClient) *AuthHandler {
	return &AuthHandler{client: c}
}

func (h *AuthHandler) RegisterLocal(c *fiber.Ctx) error {
	var body domain.Authenticate
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := h.client.Register(ctx, &pb.AuthRequest{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body domain.Authenticate
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := h.client.Login(ctx, &pb.AuthRequest{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var body domain.RefreshToken
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := h.client.RefreshToken(ctx, &pb.RefreshTokenRequest{
		RefreshToken: body.RefreshToken,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}
