package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"time"
	"wealth-vault/api-gateway/internal/domain"
	"wealth-vault/api-gateway/internal/mapper"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"
	"wealth-vault/api-gateway/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type AssetHandler struct {
	client  pb.AssetServiceClient
	storage *utils.StorageClient
}

func NewAssetHandler(c pb.AssetServiceClient, s *utils.StorageClient) *AssetHandler {
	return &AssetHandler{
		client:  c,
		storage: s,
	}
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

	if body.Amount == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Amount are required",
		})
	}

	amount, err := utils.Parseamount(body.Amount)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid amount",
		})
	}

	var pbFiles []*pb.FileInfo
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["files"]

		if len(files) > 0 {
			pbFiles = make([]*pb.FileInfo, len(files))
			var g errgroup.Group

			for i, fileHeader := range files {
				index := i
				f := fileHeader
				g.Go(func() error {
					fileData, err := f.Open()
					if err != nil {
						return err
					}
					defer fileData.Close()

					ext := filepath.Ext(f.Filename)
					newFileName := fmt.Sprintf("cash/%s-%d-%d%s", userID, time.Now().UnixNano(), index, ext)

					url, err := h.storage.UploadStream(fileData, newFileName, f.Header.Get("Content-Type"))
					if err != nil {
						return err
					}

					pbFiles[index] = &pb.FileInfo{Url: url, FileType: f.Header.Get("Content-Type")}
					return nil
				})
			}
			if err := g.Wait(); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.CreateCash(ctx, &pb.CreateCashRequest{
		Name:        body.Name,
		Amount:      amount,
		Description: body.Description,
		CreatedBy:   userID,
		Files:       pbFiles,
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

	res, err := h.client.GetCashByID(ctx, &pb.CashByIDRequest{
		UserId: userID,
		Id:     id,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	cashInfo := mapper.ToCashDomain(res.Cash)

	return c.JSON(fiber.Map{
		"status": "get cashInfo success",
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

	var newPbFiles []*pb.FileInfo

	form, err := c.MultipartForm()
	if err == nil && form != nil && len(form.File["files"]) > 0 {
		files := form.File["files"]
		newPbFiles = make([]*pb.FileInfo, len(files))
		var g errgroup.Group
		for i, fileHeader := range files {
			index := i
			f := fileHeader
			g.Go(func() (err error) {
				fileData, _ := f.Open()
				defer fileData.Close()

				newFileName := fmt.Sprintf("cash/%s-%d-%d%s", userID, time.Now().UnixNano(), index, filepath.Ext(f.Filename))
				url, err := h.storage.UploadStream(fileData, newFileName, f.Header.Get("Content-Type"))
				if err != nil {
					return err
				}

				newPbFiles[index] = &pb.FileInfo{Url: url, FileType: f.Header.Get("Content-Type")}
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Upload failed"})
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
		NewFiles:      newPbFiles,
		DeleteFileIds: req.DeleteFileIDs,
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

func (h *AssetHandler) DeleteCash(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.DeleteCash(ctx, &pb.CashByIDRequest{
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
