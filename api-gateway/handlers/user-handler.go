package handlers

import (
	"context"
	"time"

	pb "wealth-vault/api-gateway/proto/userpb"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	client pb.UserServiceClient
}

func NewUserHandler(c pb.UserServiceClient) *UserHandler {
	return &UserHandler{client: c}
}

type CreateUserDTO struct {
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Username    string `json:"username"`
	Phonenumber string `json:"phonenumber"`
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var body CreateUserDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := h.client.CreateUser(ctx, &pb.CreateUserRequest{
		Firstname:   body.Firstname,
		Lastname:    body.Lastname,
		Username:    body.Username,
		Phonenumber: body.Phonenumber,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}
