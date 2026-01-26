package handlers

import (
	"context"
	"strings"
	"time"
	"wealth-vault/api-gateway/internal/domain"
	"wealth-vault/api-gateway/internal/mapper"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"
	"wealth-vault/api-gateway/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type InvestmentHandler struct {
	client  pb.AssetServiceClient
	storage *utils.StorageClient
}

func NewInvestmentHandler(c pb.AssetServiceClient, s *utils.StorageClient) *InvestmentHandler {
	return &InvestmentHandler{
		client:  c,
		storage: s,
	}
}

func (h *InvestmentHandler) CreateInvestment(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.CreateInvestmentRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	var quantity, cost float64
	var err error
	if req.Quantity != "" {
		quantity, err = utils.Parseamount(req.Quantity)
		if err != nil {
			return err
		}
	}

	if req.CostPerPrice != "" {
		cost, err = utils.Parseamount(req.CostPerPrice)
		if err != nil {
			return err
		}
	}

	inputType := strings.ToUpper(strings.TrimSpace(req.Type))
	assetTypeEnum := pb.InvestmentType(mapper.SafeMapEnum(pb.InvestmentType_value, inputType, "INVEST_TYPE_"))
	var pbFiles []*pb.FileInfo
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["files"]

		if len(files) > 0 {
			uploadedFiles, err := utils.UploadBatchFiles(files, userID, "investment", h.storage)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			pbFiles = uploadedFiles
		}
	}

	grpcReq := &pb.CreateInvestmentRequest{
		UserId:      userID,
		Name:        req.Name,
		Symbol:      req.Symbol,
		Type:        assetTypeEnum,
		BrokerName:  req.BrokerName,
		Quantity:    quantity,
		CostPrice:   cost,
		Description: req.Description,
		NewFiles:    pbFiles,
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.CreateInvestment(ctx, grpcReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	inInfo := mapper.ToInvestDomain(res.Invest)

	return c.JSON(fiber.Map{
		"status": "create investment success",
		"data":   inInfo,
	})
}

func (h *InvestmentHandler) GetInvestment(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetInvestment(ctx, &pb.GetAssetRequest{
		UserId: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	InvestInfo := mapper.ToInvestList(res.Invest)

	return c.JSON(fiber.Map{
		"status": "get Invest Info success",
		"data":   InvestInfo,
	})
}

func (h *InvestmentHandler) GetInvestmentByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetInvestmentByID(ctx, &pb.GetAssetByIDRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	InvestInfo := mapper.ToInvestDomain(res.Invest)

	return c.JSON(fiber.Map{
		"status": "get Invest Info success",
		"data":   InvestInfo,
	})
}

func (h *InvestmentHandler) UpdateInvestment(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.UpdateInvestmentRequest)
	paths, err := utils.GetFieldMaskPaths(c, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form data",
		})
	}

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var quantity, cost float64
	if req.Quantity != "" {
		quantity, err = utils.Parseamount(req.Quantity)
		if err != nil {
			return err
		}
	}

	if req.CostPerPrice != "" {
		cost, err = utils.Parseamount(req.CostPerPrice)
		if err != nil {
			return err
		}
	}

	inputType := strings.ToUpper(strings.TrimSpace(req.Type))
	assetTypeEnum := pb.InvestmentType(mapper.SafeMapEnum(pb.InvestmentType_value, inputType, "INVEST_TYPE_"))
	var newPbFiles []*pb.FileInfo
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["files"]

		if len(files) > 0 {
			uploadedFiles, err := utils.UploadBatchFiles(files, userID, "investment", h.storage)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			newPbFiles = uploadedFiles
		}
	}

	grpcReq := &pb.UpdateInvestmentRequest{
		Id: id,
		Invest: &pb.Investment{
			UserId:      userID,
			Name:        req.Name,
			Symbol:      req.Symbol,
			Type:        assetTypeEnum,
			BrokerName:  req.BrokerName,
			Quantity:    quantity,
			CostPrice:   cost,
			Description: req.Description,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: paths,
		},
		NewFiles:      newPbFiles,
		DeleteFileIds: req.DeleteFileIDs,
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.UpdateInvestment(ctx, grpcReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	InvestInfo := mapper.ToInvestDomain(res.Invest)

	return c.JSON(fiber.Map{
		"status": "update Invest Info success",
		"data":   InvestInfo,
	})
}

func (h *InvestmentHandler) DeleteInvestment(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.DeleteInvestment(ctx, &pb.DeleteAssetRequest{
		UserId: userID,
		Id:     id,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "delete success",
		"data":   res,
	})
}
