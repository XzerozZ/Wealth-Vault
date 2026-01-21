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

func (h *AssetHandler) CreateAsset(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.CreateAssetRequest)
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
	var assetTypeEnum pb.AssetType = pb.AssetType(mapper.SafeMapEnum(pb.AssetType_value, inputType, "ASSET_TYPE_"))
	folderName := utils.GetFolderName(assetTypeEnum)
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

	grpcReq := &pb.CreateAssetRequest{
		UserId:      userID,
		Type:        assetTypeEnum,
		Name:        req.Name,
		Amount:      amount,
		Description: req.Description,
		NewFiles:    pbFiles,
	}

	detailWrapper := mapper.BuildCreateDetail(c, assetTypeEnum)
	if detailWrapper != nil {
		switch v := detailWrapper.(type) {
		case *pb.CreateAssetRequest_BankDetail:
			grpcReq.Detail = v
			grpcReq.IsIncludedInNetWorth = true
		case *pb.CreateAssetRequest_InvestmentDetail:
			grpcReq.Detail = v
			grpcReq.IsIncludedInNetWorth = true
		case *pb.CreateAssetRequest_RealEstateDetail:
			grpcReq.Detail = v
			grpcReq.IsIncludedInNetWorth = true
		case *pb.CreateAssetRequest_InsuranceDetail:
			grpcReq.Detail = v
			grpcReq.IsIncludedInNetWorth = false
		}
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.CreateAsset(ctx, grpcReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "create asset success",
		"data":   res,
	})
}

func (h *AssetHandler) GetAsset(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetAsset(ctx, &pb.GetAssetRequest{
		UserId: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	assetInfo := mapper.ToAssetList(res.Asset)

	return c.JSON(fiber.Map{
		"status": "get Asset Info success",
		"data":   assetInfo,
	})
}

func (h *AssetHandler) GetAssetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetAssetByID(ctx, &pb.GetAssetByIDRequest{
		UserId: userID,
		Id:     id,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	assetInfo := mapper.ToAssetDomain(res.Asset)

	return c.JSON(fiber.Map{
		"status": "get Asset Info success",
		"data":   assetInfo,
	})
}

func (h *AssetHandler) UpdateAsset(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.UpdateAssetRequest)
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
	var assetTypeEnum pb.AssetType = pb.AssetType(mapper.SafeMapEnum(pb.AssetType_value, inputType, "ASSET_TYPE_"))
	var amount float64
	if req.Amount != "" {
		amount, err = utils.Parseamount(req.Amount)
		if err != nil {
			return err
		}
	}

	var newPbFiles []*pb.FileInfo
	form, err := c.MultipartForm()
	folderName := utils.GetFolderName(assetTypeEnum)
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

	grpcReq := &pb.UpdateAssetRequest{
		Id:          id,
		UserId:      userID,
		Name:        req.Name,
		Amount:      amount,
		Description: req.Description,
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: paths,
		},
		NewFiles:      newPbFiles,
		DeleteFileIds: req.DeleteFileIDs,
	}

	detailWrapper := mapper.BuildUpdateDetail(c, assetTypeEnum)
	if detailWrapper != nil {
		switch v := detailWrapper.(type) {
		case *pb.UpdateAssetRequest_BankDetail:
			grpcReq.Detail = v

			paths = append(paths, "detail")
			grpcReq.UpdateMask.Paths = utils.Unique(paths)

		case *pb.UpdateAssetRequest_InvestmentDetail:
			grpcReq.Detail = v

			paths = append(paths, "detail")
			grpcReq.UpdateMask.Paths = utils.Unique(paths)

		case *pb.UpdateAssetRequest_RealEstateDetail:
			grpcReq.Detail = v

			paths = append(paths, "detail")
			grpcReq.UpdateMask.Paths = utils.Unique(paths)

		case *pb.UpdateAssetRequest_InsuranceDetail:
			grpcReq.Detail = v

			paths = append(paths, "detail")
			grpcReq.UpdateMask.Paths = utils.Unique(paths)
		}

	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.UpdateAsset(ctx, grpcReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	assetInfo := mapper.ToAssetDomain(res.Asset)

	return c.JSON(fiber.Map{
		"status": "update success",
		"data":   assetInfo,
	})
}

func (h *AssetHandler) DeleteAsset(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.DeleteAsset(ctx, &pb.DeleteAssetRequest{
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
