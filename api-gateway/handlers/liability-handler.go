package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"wealth-vault/api-gateway/internal/domain"
	"wealth-vault/api-gateway/internal/mapper"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"
	"wealth-vault/api-gateway/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type LiabilityHandler struct {
	client  pb.AssetServiceClient
	storage *utils.StorageClient
}

func NewLiabilityHandler(c pb.AssetServiceClient, s *utils.StorageClient) *LiabilityHandler {
	return &LiabilityHandler{
		client:  c,
		storage: s,
	}
}

func (h *LiabilityHandler) CreateLiability(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.CreateLiabilityRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	var principal float64
	var err error
	if req.Principal != "" {
		principal, err = utils.Parseamount(req.Principal)
		if err != nil {
			return err
		}
	}

	var interest float64
	if req.InterestRate != "" {
		interest, err = utils.Parseamount(req.InterestRate)
		if err != nil {
			return err
		}
	}

	var startAtTime, endAtTime *time.Time
	if req.StartAt != "" {
		t, err := time.Parse("2006-01-02", req.StartAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start_at format (use YYYY-MM-DD)"})
		}
		startAtTime = &t
	}

	if req.EndAt != "" {
		t, err := time.Parse("2006-01-02", req.EndAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid end_at format (use YYYY-MM-DD)"})
		}
		endAtTime = &t
	}

	inputType := strings.ToUpper(strings.TrimSpace(req.Type))
	var assetTypeEnum pb.LiabilityType = pb.LiabilityType(mapper.SafeMapEnum(pb.LiabilityType_value, inputType, "LIABILITY_TYPE_"))
	folderName := utils.GetFolderLiaName(assetTypeEnum)
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
					newFileName := fmt.Sprintf("%s/%s-%d-%d%s", folderName, userID, time.Now().UnixNano(), index, ext)

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

	grpcReq := &pb.CreateLiabilityRequest{
		UserId:       userID,
		Type:         assetTypeEnum,
		Name:         req.Name,
		Creditor:     req.Creditor,
		Principal:    principal,
		InterestRate: interest,
		Description:  req.Description,
		StartAt:      utils.ToProtoTime(startAtTime),
		EndAt:        utils.ToProtoTime(endAtTime),
		NewFiles:     pbFiles,
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.CreateLiability(ctx, grpcReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "create liability success",
		"data":   res,
	})
}

func (h *LiabilityHandler) GetLiability(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetLiability(ctx, &pb.GetLiabilityRequest{
		UserId: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	LiabilityInfo := mapper.ToLiabilityList(res.Liability)

	return c.JSON(fiber.Map{
		"status": "get Liability Info success",
		"data":   LiabilityInfo,
	})
}

func (h *LiabilityHandler) GetLiabilityByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetLiabilityByID(ctx, &pb.GetLiabilityByIDRequest{
		UserId: userID,
		Id:     id,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	liabilityInfo := mapper.ToLiabilityDomain(res.Liability)

	return c.JSON(fiber.Map{
		"status": "get Liability Info success",
		"data":   liabilityInfo,
	})
}

func (h *LiabilityHandler) UpdateLiability(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.UpdateLiabilityRequest)
	paths, err := utils.GetFieldMaskPaths(c, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form data",
		})
	}

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	inputType := strings.ToUpper(strings.TrimSpace(req.Type))
	var assetTypeEnum pb.LiabilityType = pb.LiabilityType(mapper.SafeMapEnum(pb.LiabilityType_value, inputType, "LIABILITY_TYPE_"))
	var principal float64
	if req.Principal != "" {
		principal, err = utils.Parseamount(req.Principal)
		if err != nil {
			return err
		}
	}

	var interest float64
	if req.InterestRate != "" {
		interest, err = utils.Parseamount(req.InterestRate)
		if err != nil {
			return err
		}
	}

	var startAtTime, endAtTime *time.Time
	if req.StartAt != "" {
		t, err := time.Parse("2006-01-02", req.StartAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start_at format (use YYYY-MM-DD)"})
		}
		startAtTime = &t
	}

	if req.EndAt != "" {
		t, err := time.Parse("2006-01-02", req.EndAt)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid end_at format (use YYYY-MM-DD)"})
		}
		endAtTime = &t
	}

	var newPbFiles []*pb.FileInfo
	form, err := c.MultipartForm()
	folderName := utils.GetFolderLiaName(assetTypeEnum)
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

				ext := filepath.Ext(f.Filename)
				newFileName := fmt.Sprintf("%s/%s-%d-%d%s", folderName, userID, time.Now().UnixNano(), index, ext)
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

	grpcReq := &pb.UpdateLiabilityRequest{
		Id:           id,
		UserId:       userID,
		Name:         req.Name,
		Creditor:     req.Creditor,
		Principal:    principal,
		InterestRate: interest,
		Description:  req.Description,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: paths,
		},
		StartAt:       utils.ToProtoTime(startAtTime),
		EndAt:         utils.ToProtoTime(endAtTime),
		NewFiles:      newPbFiles,
		DeleteFileIds: req.DeleteFileIDs,
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.UpdateLiability(ctx, grpcReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	liabilityInfo := mapper.ToLiabilityDomain(res.Liability)

	return c.JSON(fiber.Map{
		"status": "update success",
		"data":   liabilityInfo,
	})
}

func (h *LiabilityHandler) DeleteLiability(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.DeleteLiability(ctx, &pb.DeleteLiabilityRequest{
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
