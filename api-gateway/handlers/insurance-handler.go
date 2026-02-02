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

type InsuranceHandler struct {
	client  pb.AssetServiceClient
	storage *utils.StorageClient
}

func NewInsuranceHandler(c pb.AssetServiceClient, s *utils.StorageClient) *InsuranceHandler {
	return &InsuranceHandler{
		client:  c,
		storage: s,
	}
}

func (h *InsuranceHandler) CreateInsurance(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.CreateInsuranceRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	var conDate *time.Time
	if req.ConDate != "" {
		t, err := time.Parse("2006-01-02", req.ConDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start_at format (use YYYY-MM-DD)"})
		}
		conDate = &t
	}

	var expDate *time.Time
	if req.ExpDate != "" {
		t, err := time.Parse("2006-01-02", req.ExpDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start_at format (use YYYY-MM-DD)"})
		}
		expDate = &t
	}

	var period, amount float64
	var err error
	if req.CoveragePeriod != "" {
		period, err = utils.Parseamount(req.CoveragePeriod)
		if err != nil {
			return err
		}
	}

	if req.CoverageAmount != "" {
		amount, err = utils.Parseamount(req.CoverageAmount)
		if err != nil {
			return err
		}
	}

	inputType := strings.ToUpper(strings.TrimSpace(req.Type))
	assetTypeEnum := pb.InsuranceType(mapper.SafeMapEnum(pb.InsuranceType_value, inputType, "INSURANCE_TYPE_"))
	var pbFiles []*pb.FileInfo
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["files"]

		if len(files) > 0 {
			uploadedFiles, err := utils.UploadBatchFiles(files, userID, "insurance", h.storage)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			pbFiles = uploadedFiles
		}
	}

	grpcReq := &pb.CreateInsuranceRequest{
		UserId:         userID,
		Name:           req.Name,
		PolNum:         req.PolicyNumber,
		CompanyName:    req.CompanyName,
		CoveragePeriod: period,
		CoverageAmount: amount,
		ConDate:        utils.ToProtoTime(conDate),
		ExpDate:        utils.ToProtoTime(expDate),
		Type:           assetTypeEnum,
		Description:    req.Description,
		NewFiles:       pbFiles,
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.CreateInsurance(ctx, grpcReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	insInfo := mapper.ToInsuranceDomain(res.Insurance)

	return c.JSON(fiber.Map{
		"status": "create insurance success",
		"data":   insInfo,
	})
}

func (h *InsuranceHandler) GetInsurance(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetInsurance(ctx, &pb.GetAssetRequest{
		UserId: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	InsuranceInfo := mapper.ToInsuranceList(res.Insurance)

	return c.JSON(fiber.Map{
		"status": "get Insurance Info success",
		"data":   InsuranceInfo,
	})
}

func (h *InsuranceHandler) GetInsuranceByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetInsuranceByID(ctx, &pb.GetAssetByIDRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	InsuranceInfo := mapper.ToInsuranceDomain(res.Insurance)

	return c.JSON(fiber.Map{
		"status": "get Insurance Info success",
		"data":   InsuranceInfo,
	})
}

func (h *InsuranceHandler) UpdateInsurance(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.UpdateInsuranceRequest)
	paths, err := utils.GetFieldMaskPaths(c, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form data",
		})
	}

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var conDate *time.Time
	if req.ConDate != "" {
		t, err := time.Parse("2006-01-02", req.ConDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start_at format (use YYYY-MM-DD)"})
		}
		conDate = &t
	}

	var expDate *time.Time
	if req.ExpDate != "" {
		t, err := time.Parse("2006-01-02", req.ExpDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start_at format (use YYYY-MM-DD)"})
		}
		expDate = &t
	}

	var period, amount float64
	if req.CoveragePeriod != "" {
		period, err = utils.Parseamount(req.CoveragePeriod)
		if err != nil {
			return err
		}
	}

	if req.CoverageAmount != "" {
		amount, err = utils.Parseamount(req.CoverageAmount)
		if err != nil {
			return err
		}
	}

	inputType := strings.ToUpper(strings.TrimSpace(req.Type))
	assetTypeEnum := pb.InsuranceType(mapper.SafeMapEnum(pb.InsuranceType_value, inputType, "INSURANCE_TYPE_"))
	var newPbFiles []*pb.FileInfo
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["files"]

		if len(files) > 0 {
			uploadedFiles, err := utils.UploadBatchFiles(files, userID, "insurance", h.storage)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			newPbFiles = uploadedFiles
		}
	}

	grpcReq := &pb.UpdateInsuranceRequest{
		Id: id,
		Insurance: &pb.Insurance{
			UserId:         userID,
			Name:           req.Name,
			PolNum:         req.PolicyNumber,
			CompanyName:    req.CompanyName,
			CoveragePeriod: period,
			CoverageAmount: amount,
			ConDate:        utils.ToProtoTime(conDate),
			ExpDate:        utils.ToProtoTime(expDate),
			Type:           assetTypeEnum,
			Description:    req.Description,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: paths,
		},
		NewFiles:      newPbFiles,
		DeleteFileIds: req.DeleteFileIDs,
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.UpdateInsurance(ctx, grpcReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	InsuranceInfo := mapper.ToInsuranceDomain(res.Insurance)

	return c.JSON(fiber.Map{
		"status": "update Insurance Info success",
		"data":   InsuranceInfo,
	})
}

func (h *InsuranceHandler) DeleteInsurance(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.DeleteInsurance(ctx, &pb.DeleteAssetRequest{
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
