package middleware

import (
	"errors"
	"strings"
	"wealth-vault/api-gateway/configs"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(config configs.JWT) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing_token",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(config.Secret), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) || strings.Contains(err.Error(), "token is expired") {
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "token_expired", // Refresh
				})
			}

			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid_token",
			})
		}

		if !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid_token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid_claims",
			})
		}

		ctx.Locals("user_id", claims["user_id"])
		return ctx.Next()
	}
}
