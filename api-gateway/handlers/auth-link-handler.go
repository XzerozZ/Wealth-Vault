package handlers

import (
	"encoding/json"
	"log"
	"time"
	"wealth-vault/api-gateway/internal/domain"

	pb "wealth-vault/api-gateway/pkg/pb/proto/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func (h *AuthHandler) LinkLineAccount(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	var req domain.OAuth
	if err := c.BodyParser(&req); err != nil || req.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "access_token is required",
		})
	}

	profileReq := fasthttp.AcquireRequest()
	profileRes := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(profileReq)
	defer fasthttp.ReleaseResponse(profileRes)

	profileReq.SetRequestURI("https://api.line.me/v2/profile")
	profileReq.Header.SetMethod(fiber.MethodGet)
	profileReq.Header.Set("Authorization", "Bearer "+req.Token)

	if err := fasthttp.DoTimeout(profileReq, profileRes, 5*time.Second); err != nil {
		log.Printf("LINE Profile Error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to contact LINE server",
		})
	}

	if profileRes.StatusCode() != fiber.StatusOK {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid line access token",
		})
	}

	var profile domain.LineProfile
	if err := json.Unmarshal(profileRes.Body(), &profile); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "invalid response from LINE",
		})
	}

	if profile.UserID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid line access token",
		})
	}

	res, err := h.client.LinkLineAccount(c.Context(), &pb.LinkLineAccountRequest{
		UserId:     userID,
		LineUserId: profile.UserID,
	})

	if err != nil || !res.Success {
		log.Printf("LinkLineAccount failed user=%s line=%s err=%v", userID, profile.UserID, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to link account or already linked",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "LINE account linked successfully",
		"line_id": profile.UserID,
	})
}
