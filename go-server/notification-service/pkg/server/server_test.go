package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/notification-service/pkg/config"
)

func TestHealthHandlerReturnsStructuredStatus(t *testing.T) {
	s := &Server{cfg: &config.Config{Address: ":8080"}}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	s.healthHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("expected valid json response: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %v", body["status"])
	}

	if body["service"] != "attention-orchestration" {
		t.Fatalf("expected service attention-orchestration, got %v", body["service"])
	}
}

func TestLoadConfigRejectsInvalidPort(t *testing.T) {
	t.Setenv("PORT", "abc")

	_, err := config.LoadConfig()
	if err == nil {
		t.Fatal("expected invalid port to return an error")
	}
}

func TestAuthMiddlewareRejectsMissingApiKeyWhenConfigured(t *testing.T) {
	s := &Server{cfg: &config.Config{Address: ":8080", APIKey: "secret"}}
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	handler := s.authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestCorsMiddlewareAddsBrowserHeaders(t *testing.T) {
	s := &Server{cfg: &config.Config{Address: ":8080"}}
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	handler := s.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	handler.ServeHTTP(recorder, request)

	if recorder.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("expected cors origin header to be set, got %q", recorder.Header().Get("Access-Control-Allow-Origin"))
	}
}
