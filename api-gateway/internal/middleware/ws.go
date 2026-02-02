package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func TokenFromQuery(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return fiber.ErrUnauthorized
	}
	c.Request().Header.Set("Authorization", "Bearer "+token)
	return c.Next()
}
