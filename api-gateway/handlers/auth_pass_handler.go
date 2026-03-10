package handlers

import (
	"context"
	"time"
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/auth"

	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var body domain.ForgetPassword
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.ForgotPassword(ctx, &pb.ForgotPasswordRequest{
		Email: body.Email,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "send otp success",
		"data":   res,
	})
}

func (h *AuthHandler) VerifyForgotPasswordOTP(c *fiber.Ctx) error {
	var body domain.OTP
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.VerifyForgotPasswordOTP(ctx, &pb.VerifyOTPRequest{
		Email: body.Email,
		Otp:   body.OTP,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   res,
	})
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var body domain.ResetPassword
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if body.ResetToken == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token and Password are required",
		})
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
	defer cancel()

	res, err := h.client.ResetPassword(ctx, &pb.ResetPasswordRequest{
		ResetToken:  body.ResetToken,
		NewPassword: body.Password,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "reset password success",
		"data":   res,
	})
}
