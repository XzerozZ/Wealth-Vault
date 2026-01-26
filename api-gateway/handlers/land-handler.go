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

type LandHandler struct {
	client  pb.AssetServiceClient
	storage *utils.StorageClient
}

func NewLandHandler(c pb.AssetServiceClient, s *utils.StorageClient) *LandHandler {
	return &LandHandler{
		client:  c,
		storage: s,
	}
}

func (h *LandHandler) CreateLand(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.CreateLandRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	var area, amount float64
	var err error
	if req.Area != "" {
		area, err = utils.Parseamount(req.Area)
		if err != nil {
			return err
		}
	}

	if req.Amount != "" {
		amount, err = utils.Parseamount(req.Amount)
		if err != nil {
			return err
		}
	}

	var pbFiles []*pb.FileInfo
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["files"]

		if len(files) > 0 {
			uploadedFiles, err := utils.UploadBatchFiles(files, userID, "land", h.storage)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			pbFiles = uploadedFiles
		}
	}

	grpcReq := &pb.CreateLandRequest{
		UserId:      userID,
		Name:        req.Name,
		DeedNum:     req.DeedNum,
		Area:        area,
		Amount:      amount,
		Description: req.Description,
		Location: &pb.Location{
			Address:     req.Location.Address,
			Subdistrict: req.Location.Subdistrict,
			District:    req.Location.District,
			Province:    req.Location.Province,
			PostalCode:  req.Location.PostalCode,
		},
		NewFiles:    pbFiles,
		BuildingIds: req.ReferenceIDs,
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.CreateLand(ctx, grpcReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	LandInfo := mapper.ToLandDomain(res.Land)

	return c.JSON(fiber.Map{
		"status": "create land success",
		"data":   LandInfo,
	})
}

func (h *LandHandler) GetLand(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetLand(ctx, &pb.GetAssetRequest{
		UserId: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	LandInfo := mapper.ToLandList(res.Land)

	return c.JSON(fiber.Map{
		"status": "get Land Info success",
		"data":   LandInfo,
	})
}

func (h *LandHandler) GetLandByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetLandByID(ctx, &pb.GetAssetByIDRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	LandInfo := mapper.ToLandDomain(res.Land)

	return c.JSON(fiber.Map{
		"status": "get Land Info success",
		"data":   LandInfo,
	})
}

func (h *LandHandler) UpdateLand(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.UpdateLandRequest)
	paths, err := utils.GetFieldMaskPaths(c, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form data",
		})
	}

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var area, amount float64
	if req.Area != "" {
		area, err = utils.Parseamount(req.Area)
		if err != nil {
			return err
		}
	}

	if req.Amount != "" {
		amount, err = utils.Parseamount(req.Amount)
		if err != nil {
			return err
		}
	}

	var newPbFiles []*pb.FileInfo
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["files"]

		if len(files) > 0 {
			uploadedFiles, err := utils.UploadBatchFiles(files, userID, "land", h.storage)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			newPbFiles = uploadedFiles
		}
	}

	grpcReq := &pb.UpdateLandRequest{
		Id: id,
		Land: &pb.Land{
			UserId:      userID,
			Name:        req.Name,
			DeedNum:     req.DeedNum,
			Area:        area,
			Amount:      amount,
			Description: req.Description,
			Location: &pb.Location{
				Address:     req.Location.Address,
				Subdistrict: req.Location.Subdistrict,
				District:    req.Location.District,
				Province:    req.Location.Province,
				PostalCode:  req.Location.PostalCode,
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: paths,
		},
		NewFiles:          newPbFiles,
		BuildingIds:       req.ReferenceIDs,
		DeleteBuildingIds: req.DeleteReferenceIDs,
		DeleteFileIds:     req.DeleteFileIDs,
	}

	fmt.Println(paths)
	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.UpdateLand(ctx, grpcReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	LandInfo := mapper.ToLandDomain(res.Land)

	return c.JSON(fiber.Map{
		"status": "update Land Info success",
		"data":   LandInfo,
	})
}

func (h *LandHandler) DeleteLand(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.DeleteLand(ctx, &pb.DeleteAssetRequest{
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
