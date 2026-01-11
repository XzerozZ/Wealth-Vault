package handlers

import (
	"context"
	"time"

	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/user"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type UserHandler struct {
	client pb.UserServiceClient
}

func NewUserHandler(c pb.UserServiceClient) *UserHandler {
	return &UserHandler{client: c}
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req domain.UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var paths []string
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form-data"})
	}

	for key := range form.Value {
		paths = append(paths, key)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := h.client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id: userID,
		User: &pb.User{
			Firstname:   req.Firstname,
			Lastname:    req.Lastname,
			Username:    req.Username,
			Profile:     req.Profile,
			Phonenumber: req.Phonenumber,
			Birthday:    req.Birthday,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: paths,
		},
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "update success",
		"data":   res.User,
	})
}
