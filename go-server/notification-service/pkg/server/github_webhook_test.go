package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/notification-service/pkg/config"
	"github.com/example/notification-service/pkg/event"
)

func TestNormalizeGitHubEvent(t *testing.T) {
	s := &Server{}
	payload := map[string]interface{}{
		"repository":   map[string]interface{}{"full_name": "octo/repo"},
		"pull_request": map[string]interface{}{"number": 42},
		"sender":       map[string]interface{}{"login": "octocat"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	e, err := s.normalizeGitHubEvent("pull_request", body)
	if err != nil {
		t.Fatalf("normalize event: %v", err)
	}

	if e.Source != "github" {
		t.Fatalf("expected source github, got %s", e.Source)
	}
	if e.EventType != "pull_request" {
		t.Fatalf("expected event type pull_request, got %s", e.EventType)
	}
	if e.Repository != "octo/repo" {
		t.Fatalf("expected repository octo/repo, got %s", e.Repository)
	}
}

func TestGitHubWebhookHandlerPersistsNormalizedEvent(t *testing.T) {
	repo := event.NewInMemoryRepository()
	s := &Server{cfg: &config.Config{Address: ":8080"}, eventRepo: repo}

	payload := []byte(`{"repository":{"full_name":"octo/repo"},"pull_request":{"number":42},"sender":{"login":"octocat"}}`)
	req := httptest.NewRequest(http.MethodPost, "/webhooks/github", bytes.NewReader(payload))
	req.Header.Set("X-GitHub-Event", "pull_request")
	rr := httptest.NewRecorder()

	s.githubWebhookHandler(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, rr.Code)
	}

	stored, err := repo.List()
	if err != nil {
		t.Fatalf("list events: %v", err)
	}
	if len(stored) != 1 {
		t.Fatalf("expected 1 stored event, got %d", len(stored))
	}
}
