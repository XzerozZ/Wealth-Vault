package domain

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestShareJSONContract(t *testing.T) {
	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		var req ShareItemRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}
		if len(req.ItemIDs) != 1 || req.ItemIDs[0] != "item-id" || len(req.ItemTypes) != 1 || req.ItemTypes[0] != "cash" {
			return fiber.ErrBadRequest
		}
		if len(req.Friends) != 1 || req.Friends[0].ID != "friend-id" {
			return fiber.ErrBadRequest
		}
		return c.SendStatus(200)
	})
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"item_ids":["item-id"],"item_types":["cash"],"friends":[{"id":"friend-id"}]}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("status=%d", res.StatusCode)
	}
}
