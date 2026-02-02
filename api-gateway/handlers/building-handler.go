package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"
	"wealth-vault/api-gateway/internal/domain"
	"wealth-vault/api-gateway/internal/mapper"
	pb "wealth-vault/api-gateway/pkg/pb/proto/asset"
	"wealth-vault/api-gateway/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type BuildingHandler struct {
	client  pb.AssetServiceClient
	storage *utils.StorageClient
}

func NewBuildingHandler(c pb.AssetServiceClient, s *utils.StorageClient) *BuildingHandler {
	return &BuildingHandler{
		client:  c,
		storage: s,
	}
}

func (h *BuildingHandler) CreateBuilding(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.CreateBuildingRequest)
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

	inputType := strings.ToUpper(strings.TrimSpace(req.Type))
	assetTypeEnum := pb.BuildingType(mapper.SafeMapEnum(pb.BuildingType_value, inputType, "BUILDING_TYPE_"))
	var pbFiles []*pb.FileInfo
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["files"]

		if len(files) > 0 {
			uploadedFiles, err := utils.UploadBatchFiles(files, userID, "building", h.storage)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			pbFiles = uploadedFiles
		}
	}

	grpcReq := &pb.CreateBuildingRequest{
		UserId:      userID,
		Name:        req.Name,
		Type:        assetTypeEnum,
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
		NewFiles: pbFiles,
		LandIds:  req.ReferenceIDs,
		InsIds:   req.InsIDs,
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.CreateBuilding(ctx, grpcReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	buInfo := mapper.ToBuildingDomain(res.Building)

	return c.JSON(fiber.Map{
		"status": "create building success",
		"data":   buInfo,
	})
}

func (h *BuildingHandler) GetBuilding(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetBuilding(ctx, &pb.GetAssetRequest{
		UserId: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	BuildingInfo := mapper.ToBuildingList(res.Building)

	return c.JSON(fiber.Map{
		"status": "get Building Info success",
		"data":   BuildingInfo,
	})
}

func (h *BuildingHandler) GetBuildingByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.GetBuildingByID(ctx, &pb.GetAssetByIDRequest{
		Id:     id,
		UserId: userID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	BuildingInfo := mapper.ToBuildingDomain(res.Building)

	return c.JSON(fiber.Map{
		"status": "get Building Info success",
		"data":   BuildingInfo,
	})
}

func (h *BuildingHandler) UpdateBuilding(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	req := new(domain.UpdateBuildingRequest)
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

	inputType := strings.ToUpper(strings.TrimSpace(req.Type))
	assetTypeEnum := pb.BuildingType(mapper.SafeMapEnum(pb.BuildingType_value, inputType, "BUILDING_TYPE_"))
	var newPbFiles []*pb.FileInfo
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		files := form.File["files"]

		if len(files) > 0 {
			uploadedFiles, err := utils.UploadBatchFiles(files, userID, "building", h.storage)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			newPbFiles = uploadedFiles
		}
	}

	grpcReq := &pb.UpdateBuildingRequest{
		Id: id,
		Building: &pb.Building{
			UserId:      userID,
			Name:        req.Name,
			Type:        assetTypeEnum,
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
		NewFiles:      newPbFiles,
		LandIds:       req.ReferenceIDs,
		DeleteLandIds: req.DeleteReferenceIDs,
		InsIds:        req.InsIDs,
		DeleteInsIds:  req.DeleteInsIDs,
		DeleteFileIds: req.DeleteFileIDs,
	}

	fmt.Println(paths)
	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.UpdateBuilding(ctx, grpcReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	BuildingInfo := mapper.ToBuildingDomain(res.Building)

	return c.JSON(fiber.Map{
		"status": "update Building Info success",
		"data":   BuildingInfo,
	})
}

func (h *BuildingHandler) DeleteBuilding(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.DeleteBuilding(ctx, &pb.DeleteAssetRequest{
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
