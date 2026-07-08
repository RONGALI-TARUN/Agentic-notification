package learning

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSummarizeReturnsModelResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"reduce noise"}`))
	}))
	defer server.Close()

	service := NewService(server.URL)
	result, err := service.Summarize("user dismissed alerts")
	if err != nil {
		t.Fatalf("summarize feedback: %v", err)
	}
	if result != "reduce noise" {
		t.Fatalf("expected reduce noise, got %s", result)
	}
}
