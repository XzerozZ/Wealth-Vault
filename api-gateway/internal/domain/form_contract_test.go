package domain

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func decodeMultipart(t *testing.T, fields map[string]string, output interface{}) {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		if err := c.BodyParser(output); err != nil {
			return err
		}
		return c.JSON(output)
	})
	req := httptest.NewRequest("POST", "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("status = %d", res.StatusCode)
	}
}

func TestAssetUploadFieldContract(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	file, err := writer.CreateFormFile("files", "sample.pdf")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.Write([]byte("test-file")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}
		if len(form.File["files"]) != 1 {
			return fiber.ErrBadRequest
		}
		return c.SendStatus(200)
	})
	req := httptest.NewRequest("POST", "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("status=%d", res.StatusCode)
	}
}

func TestAssetMultipartContract(t *testing.T) {
	cases := []struct {
		name   string
		fields map[string]string
		output interface{}
		expect map[string]string
	}{
		{"account", map[string]string{"name": "Bank", "bank_name": "Test", "bank_account": "123", "type": "savings", "amount": "100"}, &CreateAccountRequest{}, map[string]string{"name": "Bank", "bank_account": "123", "amount": "100"}},
		{"cash", map[string]string{"name": "Cash", "amount": "100"}, &CreateCashRequest{}, map[string]string{"name": "Cash", "amount": "100"}},
		{"invest", map[string]string{"name": "Stock", "symbol": "AAPL", "quantity": "2", "cost_per_price": "100"}, &CreateInvestmentRequest{}, map[string]string{"name": "Stock", "symbol": "AAPL", "quantity": "2", "cost_per_price": "100"}},
		{"insurance", map[string]string{"name": "Policy", "policy_number": "ABC", "coverage_amount": "1000"}, &CreateInsuranceRequest{}, map[string]string{"name": "Policy", "policy_number": "ABC", "coverage_amount": "1000"}},
		{"land", map[string]string{"name": "Land", "deed_num": "1", "location.address": "Bangkok", "location.province": "Bangkok"}, &CreateLandRequest{}, map[string]string{"name": "Land", "deed_num": "1"}},
		{"building", map[string]string{"name": "Home", "location.address": "Bangkok", "location.province": "Bangkok"}, &CreateBuildingRequest{}, map[string]string{"name": "Home"}},
		{"liability", map[string]string{"name": "Loan", "principal": "1000", "interest_rate": "3"}, &CreateLiabilityRequest{}, map[string]string{"name": "Loan", "principal": "1000", "interest_rate": "3"}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			decodeMultipart(t, test.fields, test.output)
			encoded, err := json.Marshal(test.output)
			if err != nil {
				t.Fatal(err)
			}
			var parsed map[string]any
			if err := json.Unmarshal(encoded, &parsed); err != nil {
				t.Fatal(err)
			}
			for key, want := range test.expect {
				if parsed[key] != want {
					t.Fatalf("%s = %v, want %q", key, parsed[key], want)
				}
			}
			if test.name == "land" {
				got := test.output.(*CreateLandRequest)
				if got.Location.Address != "Bangkok" || got.Location.Province != "Bangkok" {
					t.Fatalf("location not parsed: %s", encoded)
				}
			}
			if test.name == "building" {
				got := test.output.(*CreateBuildingRequest)
				if got.Location.Address != "Bangkok" || got.Location.Province != "Bangkok" {
					t.Fatalf("location not parsed: %s", encoded)
				}
			}
			if string(encoded) == "{}" {
				t.Fatalf("empty payload: %s", encoded)
			}
		})
	}
}
