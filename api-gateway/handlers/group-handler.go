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

type GroupHandler struct {
	client  pb.UserServiceClient
	storage *utils.StorageClient
}

func NewGroupHandler(c pb.UserServiceClient, s *utils.StorageClient) *GroupHandler {
	return &GroupHandler{
		client:  c,
		storage: s,
	}
}

func (h *GroupHandler) CreateGroup(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.CreateGroupRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
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

		req.GroupProfile = url
	}

	grpcReq := &pb.CreateGroupRequest{
		Name:      req.GroupName,
		Profile:   req.GroupProfile,
		CreatorId: userID,
		MemberIds: req.MemberIDS,
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.CreateGroup(ctx, grpcReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	groupInfo := mapper.ToGroupDomain(res.Group)

	return c.JSON(fiber.Map{
		"status": "create group success",
		"data":   groupInfo,
	})
}

func (h *GroupHandler) GetMember(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetMember(ctx, &pb.GetGroupMembersRequest{
		GroupId: id,
		UserId:  userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	members := mapper.ToUserList(res.Members)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"members": members,
			"total":   res.Total,
		},
	})
}

func (h *GroupHandler) AllGetGroup(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetAllGroup(ctx, &pb.AllGroupRequest{
		UserId: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	groupInfo := mapper.ToGroupList(res.Group)

	return c.JSON(fiber.Map{
		"status": "get group success",
		"data":   groupInfo,
	})
}

func (h *GroupHandler) GetGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetGroup(ctx, &pb.GetGroupRequest{
		GroupId: id,
		UserId:  userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	groupInfo := mapper.ToGroupDomain(res.Group)

	return c.JSON(fiber.Map{
		"status": "get group success",
		"data":   groupInfo,
	})
}

func (h *GroupHandler) UpdateGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req domain.UpdateGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	paths, err := utils.GetFieldMaskPaths(c, &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid form data",
		})
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

		req.GroupProfile = url
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

	res, err := h.client.UpdateGroup(ctx, &pb.UpdateGroupRequest{
		Id:      id,
		UserId:  userID,
		Name:    req.GroupName,
		Profile: req.GroupProfile,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: paths,
		},
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	groupInfo := mapper.ToGroupDomain(res.Group)

	return c.JSON(fiber.Map{
		"status": "update group success",
		"data":   groupInfo,
	})
}

func (h *GroupHandler) AddMember(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req domain.AddMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.AddGroupMember(ctx, &pb.AddMemberRequest{
		GroupId:       id,
		UserId:        userID,
		TargetUserIds: req.TargetUserIDS,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "add member(s) success",
		"data":    res.Success,
	})
}

func (h *GroupHandler) RemoveMember(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req domain.RemoveMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.RemoveMember(ctx, &pb.RemoveMemberRequest{
		GroupId:        id,
		UserId:         userID,
		TargetMemberId: req.TargetUserID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "remove member success",
		"data":    res.Success,
	})
}

func (h *GroupHandler) LeaveGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.LeaveGroup(ctx, &pb.LeaveGroupRequest{
		GroupId: id,
		UserId:  userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "leave group success",
		"data":    res.Success,
	})
}

func (h *GroupHandler) GrantAccess(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req domain.GrantAccessRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GrantGroupItemAccess(ctx, &pb.GrantAccessRequest{
		GroupId:      id,
		OwnerUserId:  userID,
		TargetUserId: req.TargetUserID,
		GroupItemIds: req.ItemIDS,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "grant access success",
		"data":    res.Success,
	})
}

func (h *GroupHandler) DeleteGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.DeleteGroup(ctx, &pb.DeleteGroupRequest{
		GroupId: id,
		UserId:  userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "delete group success",
		"data":    res.Success,
	})
}
