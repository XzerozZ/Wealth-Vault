package handlers

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
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

	var sharedAgePtr *int32
	if req.SharedAge != "" {
		val, err := strconv.Atoi(req.SharedAge)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid shared_age format"})
		}
		fVal := int32(val)
		sharedAgePtr = &fVal
	}

	var sharedEnabledPtr *bool
	if req.SharedEnabled != "" {
		val, err := strconv.ParseBool(req.SharedEnabled)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid shared_enabled format"})
		}
		sharedEnabledPtr = &val
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
		Id:            userID,
		Firstname:     req.Firstname,
		Lastname:      req.Lastname,
		Username:      req.Username,
		Profile:       req.Profile,
		Phonenumber:   req.Phonenumber,
		Birthday:      utils.ToProtoTime(birthDate),
		Sharedage:     sharedAgePtr,
		Sharedenabled: sharedEnabledPtr,
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
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req domain.AddFriendRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.AddFriend(ctx, &pb.FriendRequest{
		Id:     userID,
		UserId: req.RequesterID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "add friend success",
		"data":   res,
	})
}

func (h *UserHandler) AcceptFriend(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req domain.AcceptFriendRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.AcceptFriend(ctx, &pb.AcceptFriendRequest{
		UserId:      userID,
		RequesterId: req.RequesterID,
		Action:      req.Action,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "add friend success",
		"data":   res,
	})
}

func (h *UserHandler) GetFriendProfile(c *fiber.Ctx) error {
	targetID := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	shareRes, err := h.client.GetItemsSharedByFriend(c.Context(), &pb.GetItemsSharedByFriendRequest{
		UserId:   userID,
		FriendId: targetID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch shared items IDs: " + err.Error()})
	}

	res, err := h.client.GetUser(ctx, &pb.GetUserByIDRequest{
		Id: targetID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	userInfo := mapper.ToUserDomain(res.User)
	itemInfo := mapper.MapAllFriendItemsToDomain(shareRes.AssetDetail)

	return c.JSON(fiber.Map{
		"status": "get userInfo success",
		"data": fiber.Map{
			"user_info":    userInfo,
			"item_preview": itemInfo,
		},
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

	friendInfo := mapper.ToUserList(res.Friends)

	return c.JSON(fiber.Map{
		"status": "get friend success",
		"data": fiber.Map{
			"friends": friendInfo,
		},
	})
}

func (h *UserHandler) GetPendingRequests(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetPendingRequests(ctx, &pb.GetUserByIDRequest{
		Id: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	friendInfo := mapper.ToUserList(res.Friends)

	return c.JSON(fiber.Map{
		"status": "get friend pending request  success",
		"data": fiber.Map{
			"friends": friendInfo,
		},
	})
}

func (h *UserHandler) ToggleCloseFriend(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req domain.SetCloseFriendRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var isclose bool
	if req.IsClose != "" {
		val, err := strconv.ParseBool(req.IsClose)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid shared_enabled format"})
		}
		isclose = val
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	res, err := h.client.SetCloseFriend(ctx, &pb.SetCloseFriendRequest{
		UserId:   userID,
		FriendId: req.FriendID,
		IsClose:  isclose,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "set closed friend success",
		"data":   res,
	})
}

func (h *UserHandler) GetCloseFriends(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	res, err := h.client.GetCloseFriends(ctx, &pb.GetCloseFriendsRequest{
		UserId: userID,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	friendInfo := mapper.ToUserList(res.Friends)
	log.Println(friendInfo[0].IsClose)
	return c.JSON(fiber.Map{
		"data": friendInfo,
	})
}
