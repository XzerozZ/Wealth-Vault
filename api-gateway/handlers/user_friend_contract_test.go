package handlers

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	pb "wealth-vault/api-gateway/pkg/pb/proto/user"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

type friendClientSpy struct {
	pb.UserServiceClient
	action string
	calls  int
}

func (s *friendClientSpy) AcceptFriend(_ context.Context, req *pb.AcceptFriendRequest, _ ...grpc.CallOption) (*pb.FriendResponse, error) {
	s.calls++
	s.action = req.Action
	return &pb.FriendResponse{}, nil
}

func TestFriendActionContract(t *testing.T) {
	for _, test := range []struct {
		action, want string
		status       int
	}{
		{"reject", "DECLINE", 200},
		{"DECLINE", "DECLINE", 200},
		{"accept", "ACCEPT", 200},
		{"unexpected", "", 400},
	} {
		t.Run(test.action, func(t *testing.T) {
			spy := &friendClientSpy{}
			app := fiber.New()
			app.Post("/friend/accept", func(c *fiber.Ctx) error { c.Locals("user_id", "test-user"); return c.Next() }, NewUserHandler(spy, nil).AcceptFriend)
			req := httptest.NewRequest("POST", "/friend/accept", strings.NewReader(`{"requester_id":"friend-id","action":"`+test.action+`"}`))
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != test.status || spy.action != test.want {
				t.Fatalf("status=%d action=%q calls=%d", res.StatusCode, spy.action, spy.calls)
			}
			if test.status == 400 && spy.calls != 0 {
				t.Fatal("invalid action reached user service")
			}
		})
	}
}
