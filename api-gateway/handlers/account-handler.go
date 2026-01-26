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

type AccountHandler struct {
	client  pb.AssetServiceClient
	storage *utils.StorageClient
}

func NewAccountHandler(c pb.AssetServiceClient, s *utils.StorageClient) *AccountHandler {
	return &AccountHandler{
		client:  c,
		storage: s,
	}
}

func (h *AccountHandler) CreateAccount(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.CreateAccountRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	var amount float64
	var err error
	if req.Amount != "" {
		amount, err = utils.Parseamount(req.Amount)
		if err != nil {
			return err
		}
	}

	inputType := strings.ToUpper(strings.TrimSpace(req.Type))
	assetTypeEnum := pb.BankAccType(mapper.SafeMapEnum(pb.BankAccType_value, inputType, "BANK_ACC_TYPE_"))
	var pbFiles []*pb.FileInfo
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["files"]

		if len(files) > 0 {
			uploadedFiles, err := utils.UploadBatchFiles(files, userID, "account", h.storage)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			pbFiles = uploadedFiles
		}
	}

	grpcReq := &pb.CreateAccountRequest{
		UserId:      userID,
		Name:        req.Name,
		BankName:    req.BankName,
		BankAcc:     req.BankAccount,
		Type:        assetTypeEnum,
		Amount:      amount,
		Description: req.Description,
		NewFiles:    pbFiles,
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.CreateAccount(ctx, grpcReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	accInfo := mapper.ToAccountDomain(res.Account)

	return c.JSON(fiber.Map{
		"status": "create account success",
		"data":   accInfo,
	})
}

func (h *AccountHandler) GetAccount(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetAccount(ctx, &pb.GetAssetRequest{
		UserId: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	AccountInfo := mapper.ToAccountList(res.Account)

	return c.JSON(fiber.Map{
		"status": "get Account Info success",
		"data":   AccountInfo,
	})
}

func (h *AccountHandler) GetAccountByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetAccountByID(ctx, &pb.GetAssetByIDRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	AccountInfo := mapper.ToAccountDomain(res.Account)

	return c.JSON(fiber.Map{
		"status": "get Account Info success",
		"data":   AccountInfo,
	})
}

func (h *AccountHandler) UpdateAccount(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.UpdateAccountRequest)
	paths, err := utils.GetFieldMaskPaths(c, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form data",
		})
	}

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var amount float64
	if req.Amount != "" {
		amount, err = utils.Parseamount(req.Amount)
		if err != nil {
			return err
		}
	}

	inputType := strings.ToUpper(strings.TrimSpace(req.Type))
	assetTypeEnum := pb.BankAccType(mapper.SafeMapEnum(pb.BankAccType_value, inputType, "BANK_ACC_TYPE_"))
	var newPbFiles []*pb.FileInfo
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["files"]

		if len(files) > 0 {
			uploadedFiles, err := utils.UploadBatchFiles(files, userID, "account", h.storage)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			newPbFiles = uploadedFiles
		}
	}

	grpcReq := &pb.UpdateAccountRequest{
		Id: id,
		Acc: &pb.Account{
			UserId:      userID,
			Name:        req.Name,
			BankName:    req.BankName,
			BankAcc:     req.BankAccount,
			Amount:      amount,
			Type:        assetTypeEnum,
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

	res, err := h.client.UpdateAccount(ctx, grpcReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	AccountInfo := mapper.ToAccountDomain(res.Account)

	return c.JSON(fiber.Map{
		"status": "update Account Info success",
		"data":   AccountInfo,
	})
}

func (h *AccountHandler) DeleteAccount(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.DeleteAccount(ctx, &pb.DeleteAssetRequest{
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
