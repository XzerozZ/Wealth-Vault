package handlers

import (
	"context"
	"strings"
	"time"
	pb "wealth-vault/api-gateway/pkg/pb/proto/auth"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type lineExchangeBody struct {
	IDToken   string `json:"idToken"`
	ChannelID string `json:"channelId"`
}

type lineChallengeBody struct {
	ChallengeID string `json:"challengeId"`
	ChannelID   string `json:"channelId"`
	Action      string `json:"action"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
}

func liffError(c *fiber.Ctx, err error) error {
	code := status.Code(err)
	switch code {
	case codes.Unauthenticated, codes.PermissionDenied:
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "challenge_rejected"})
	case codes.InvalidArgument:
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"code": "invalid_request"})
	case codes.FailedPrecondition, codes.Unavailable:
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"code": "auth_unavailable"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "auth_error"})
	}
}

func liffResponse(c *fiber.Ctx, res *pb.LineAuthResponse) error {
	c.Set("Cache-Control", "no-store")
	if res.Status == "authenticated" {
		return c.JSON(fiber.Map{"status": "authenticated", "user": fiber.Map{"id": res.UserId, "displayName": res.DisplayName},
			"accessToken": res.AccessToken, "refreshToken": res.RefreshToken, "accessTokenExpiresIn": res.AccessTokenExpiresIn})
	}
	if res.Status == "link_required" || res.Status == "onboarding_required" {
		return c.JSON(fiber.Map{"status": res.Status, "challengeId": res.ChallengeId})
	}
	return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"code": "invalid_auth_response"})
}

func (h *AuthHandler) ExchangeLine(c *fiber.Ctx) error {
	if len(c.Body()) > 8192 {
		return c.SendStatus(fiber.StatusRequestEntityTooLarge)
	}
	var body lineExchangeBody
	if err := c.BodyParser(&body); err != nil || len(body.IDToken) < 20 || len(body.IDToken) > 6000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "invalid_id_token"})
	}
	ctx, cancel := context.WithTimeout(c.UserContext(), 7*time.Second)
	defer cancel()
	res, err := h.client.ExchangeLine(ctx, &pb.LineExchangeRequest{IdToken: body.IDToken, ChannelId: body.ChannelID})
	if err != nil {
		return liffError(c, err)
	}
	return liffResponse(c, res)
}

func (h *AuthHandler) CompleteLineChallenge(c *fiber.Ctx) error {
	if len(c.Body()) > 8192 {
		return c.SendStatus(fiber.StatusRequestEntityTooLarge)
	}
	var body lineChallengeBody
	if err := c.BodyParser(&body); err != nil || len(body.ChallengeID) != 64 || (body.Action != "link" && body.Action != "register") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "invalid_challenge"})
	}
	ctx, cancel := context.WithTimeout(c.UserContext(), 7*time.Second)
	defer cancel()
	res, err := h.client.CompleteLineChallenge(ctx, &pb.LineChallengeRequest{ChallengeId: body.ChallengeID,
		ChannelId: body.ChannelID, Action: body.Action, Email: strings.TrimSpace(body.Email), Password: body.Password,
		Username: body.Username, FirstName: body.FirstName, LastName: body.LastName})
	if err != nil {
		return liffError(c, err)
	}
	return liffResponse(c, res)
}

func (h *AuthHandler) RevokeLIFFSession(c *fiber.Ctx) error {
	if len(c.Body()) > 1024 {
		return c.SendStatus(fiber.StatusRequestEntityTooLarge)
	}
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.BodyParser(&body); err != nil || len(body.RefreshToken) < 20 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "invalid_refresh_token"})
	}
	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()
	_, err := h.client.RevokeSession(ctx, &pb.RevokeSessionRequest{RefreshToken: body.RefreshToken})
	if err != nil {
		return liffError(c, err)
	}
	c.Set("Cache-Control", "no-store")
	return c.JSON(fiber.Map{"success": true})
}
