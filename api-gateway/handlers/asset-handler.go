package handlers

import (
	"context"
	"fmt"
	"time"
	"wealth-vault/api-gateway/internal/domain"
	"wealth-vault/api-gateway/internal/mapper"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"
	"wealth-vault/api-gateway/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type AssetHandler struct {
	client pb.AssetServiceClient
}

func NewAssetHandler(c pb.AssetServiceClient) *AssetHandler {
	return &AssetHandler{client: c}
}

func (h *AssetHandler) CreateCash(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var body domain.CreateCash
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	amount, err := utils.Parseamount(body.Amount)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.CreateCash(ctx, &pb.CreateCashRequest{
		Name:        body.Name,
		Amount:      amount,
		Description: body.Description,
		CreatedBy:   userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "create asset success",
		"data":   res,
	})
}

func (h *AssetHandler) GetCash(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetCash(ctx, &pb.GetCashRequest{
		UserId: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	cashInfo := mapper.ToCashList(res.Cash)

	return c.JSON(fiber.Map{
		"status": "get userInfo success",
		"data":   cashInfo,
	})
}

func (h *AssetHandler) GetCashByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetCashByID(ctx, &pb.GetCashByIDRequest{
		UserId: userID,
		Id:     id,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	cashInfo := mapper.ToCashDomain(res.Cash)

	return c.JSON(fiber.Map{
		"status": "get userInfo success",
		"data":   cashInfo,
	})
}

func (h *AssetHandler) UpdateCash(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req domain.UpdateCashRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	paths, err := utils.GetFieldMaskPaths(c, &req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid form data",
		})
	}

	var amount float64
	if req.Amount != "" {
		amount, err = utils.Parseamount(req.Amount)
		if err != nil {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.UpdateCash(ctx, &pb.UpdateCashRequest{
		Id: id,
		Cash: &pb.Cash{
			Name:        req.Name,
			Amount:      amount,
			Description: req.Description,
			CreatedBy:   userID,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: paths,
		},
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	cashInfo := mapper.ToCashDomain(res.Cash)

	return c.JSON(fiber.Map{
		"status": "update success",
		"data":   cashInfo,
	})
}
