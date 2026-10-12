package line

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestSendTextMessageRetriesWithSameKey(t *testing.T) {
	requests := 0
	key := uuid.New().String()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Header.Get("X-Line-Retry-Key") != key || r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("unexpected LINE request headers")
		}
		if requests == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	client := &lineClient{channelAccessToken: "test-token", client: server.Client(), endpoint: server.URL}
	if err := client.SendTextMessage(context.Background(), "U123", "test", key); err != nil {
		t.Fatal(err)
	}
	if requests != 2 {
		t.Fatalf("expected two attempts, got %d", requests)
	}
}

func TestSendTextMessageDoesNotRetryClientError(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()
	client := &lineClient{channelAccessToken: "test-token", client: server.Client(), endpoint: server.URL}
	if err := client.SendTextMessage(context.Background(), "U123", "test", uuid.New().String()); err == nil {
		t.Fatal("expected 400 to fail")
	}
	if requests != 1 {
		t.Fatalf("expected no retry, got %d attempts", requests)
	}
}
