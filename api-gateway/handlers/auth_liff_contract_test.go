package handlers

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	pb "wealth-vault/api-gateway/pkg/pb/proto/auth"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

type liffClientSpy struct{ pb.AuthServiceClient }

func (*liffClientSpy) ExchangeLine(_ context.Context, req *pb.LineExchangeRequest, _ ...grpc.CallOption) (*pb.LineAuthResponse, error) {
	if req.IdToken == "invalid-id-token-but-long-enough" {
		return &pb.LineAuthResponse{Status: "link_required", ChallengeId: strings.Repeat("a", 64)}, nil
	}
	return &pb.LineAuthResponse{Status: "authenticated", UserId: "test-user", DisplayName: "Tester", AccessToken: "access", RefreshToken: "refresh", AccessTokenExpiresIn: 900}, nil
}

func TestLINEExchangeGatewayContract(t *testing.T) {
	app := fiber.New()
	app.Post("/", NewAuthHandler(&liffClientSpy{}).ExchangeLine)
	for _, tc := range []struct{ token, status string }{
		{"invalid-id-token-but-long-enough", "link_required"},
		{"valid-id-token-but-long-enough", "authenticated"},
	} {
		req := httptest.NewRequest("POST", "/", strings.NewReader(`{"idToken":"`+tc.token+`","channelId":"test-channel"}`))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		_ = res.Body.Close()
		if res.StatusCode != 200 || body["status"] != tc.status {
			t.Fatalf("status=%d body=%v", res.StatusCode, body)
		}
		if tc.status == "authenticated" {
			if body["accessToken"] != "access" || body["refreshToken"] != "refresh" || body["accessTokenExpiresIn"] != float64(900) {
				t.Fatalf("invalid auth payload: %v", body)
			}
		} else if body["challengeId"] != strings.Repeat("a", 64) {
			t.Fatalf("invalid challenge payload: %v", body)
		}
	}
}
