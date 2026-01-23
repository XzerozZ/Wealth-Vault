package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"wealth-vault/api-gateway/internal/domain"
	"wealth-vault/api-gateway/internal/mapper"
	pb "wealth-vault/api-gateway/pkg/pb/proto/user"
	"wealth-vault/api-gateway/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type UserHandler struct {
	client  pb.UserServiceClient
	storage *utils.StorageClient
}

func NewUserHandler(c pb.UserServiceClient, s *utils.StorageClient) *UserHandler {
	return &UserHandler{
		client:  c,
		storage: s,
	}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetUser(ctx, &pb.GetUserByIDRequest{
		Id: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	userInfo := mapper.ToUserDomain(res.User)

	return c.JSON(fiber.Map{
		"status": "get userInfo success",
		"data":   userInfo,
	})
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

	paths, err := utils.GetFieldMaskPaths(c, &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid form data",
		})
	}

	var birthDate *time.Time
	if req.Birthday != "" {
		t, err := time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start_at format (use YYYY-MM-DD)"})
		}
		birthDate = &t
	}

	fileHeader, err := c.FormFile("profile_image")
	if err == nil {

		fileData, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot open file"})
		}
		defer fileData.Close()

		ext := filepath.Ext(fileHeader.Filename)
		newFileName := fmt.Sprintf("avatars/%s-%d%s", userID, time.Now().Unix(), ext)
		contentType := fileHeader.Header.Get("Content-Type")

		url, err := h.storage.UploadStream(fileData, newFileName, contentType)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Upload failed: " + err.Error()})
		}

		req.Profile = url
		hasProfile := false
		for _, p := range paths {
			if p == "profile" {
				hasProfile = true
				break
			}
		}
		if !hasProfile {
			paths = append(paths, "profile")
		}
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:          userID,
		Firstname:   req.Firstname,
		Lastname:    req.Lastname,
		Username:    req.Username,
		Profile:     req.Profile,
		Phonenumber: req.Phonenumber,
		Birthday:    utils.ToProtoTime(birthDate),
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: paths,
		},
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	userInfo := mapper.ToUserDomain(res.User)

	return c.JSON(fiber.Map{
		"status": "update success",
		"data":   userInfo,
	})
}

func (h *UserHandler) AddFriend(c *fiber.Ctx) error {
	id := c.Params("friendID")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.AddFriend(ctx, &pb.FriendRequest{
		Id:     userID,
		UserId: id,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "add friend success",
		"data":   res,
	})
}

func (h *UserHandler) GetFriendList(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetFriendList(ctx, &pb.GetUserByIDRequest{
		Id: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	userInfo := mapper.ToUserDomain(res.User)
	friendInfo := mapper.ToUserList(res.Friends)

	return c.JSON(fiber.Map{
		"status": "get friend success",
		"data": fiber.Map{
			"user":    userInfo,
			"friends": friendInfo,
		},
	})
}
