package pipeline

import (
	"testing"

	"github.com/example/notification-service/pkg/delivery"
	"github.com/example/notification-service/pkg/event"
	"github.com/example/notification-service/pkg/notification"
	"github.com/example/notification-service/pkg/policy"
)

type recordingProvider struct {
	sent []notification.NotificationRequest
}

func (r *recordingProvider) Send(req notification.NotificationRequest) error {
	r.sent = append(r.sent, req)
	return nil
}

func (r *recordingProvider) Name() string { return "recording" }

func newTestPipeline(rules []policy.Rule) (*Pipeline, *event.InMemoryRepository, *recordingProvider) {
	repo := event.NewInMemoryRepository()
	provider := &recordingProvider{}
	p := New(
		policy.NewEngine(rules),
		delivery.NewService(provider),
		repo,
		Config{DefaultEmailRecipients: []string{"ops@example.com"}},
	)
	return p, repo, provider
}

func TestPipelineDeliversCriticalEvents(t *testing.T) {
	p, repo, provider := newTestPipeline([]policy.Rule{{
		Name: "security",
		When: policy.RuleCondition{Source: "github", EventType: "security_advisory"},
		Then: policy.RuleOutcome{AttentionClass: "critical", Reason: "security vulnerability"},
	}})

	result, err := p.Handle(event.Event{
		ID: "evt-1", Source: "github", EventType: "security_advisory", Repository: "octo/repo",
	})
	if err != nil {
		t.Fatalf("handle event: %v", err)
	}

	if result.AttentionClass != event.AttentionCritical {
		t.Fatalf("expected critical, got %s", result.AttentionClass)
	}
	if result.DecisionState.Status != event.DecisionSent {
		t.Fatalf("expected sent status, got %s", result.DecisionState.Status)
	}
	if len(provider.sent) != 1 {
		t.Fatalf("expected 1 delivery, got %d", len(provider.sent))
	}

	stored, _ := repo.List()
	if len(stored) != 1 {
		t.Fatalf("expected 1 stored event, got %d", len(stored))
	}
}

func TestPipelineQueuesStandardEvents(t *testing.T) {
	p, repo, provider := newTestPipeline([]policy.Rule{{
		Name: "default",
		When: policy.RuleCondition{Source: "github"},
		Then: policy.RuleOutcome{AttentionClass: "standard", Reason: "default routing"},
	}})

	result, err := p.Handle(event.Event{
		ID: "evt-2", Source: "github", EventType: "pull_request", Repository: "octo/repo",
	})
	if err != nil {
		t.Fatalf("handle event: %v", err)
	}

	if result.DecisionState.Status != event.DecisionQueued {
		t.Fatalf("expected queued status, got %s", result.DecisionState.Status)
	}
	if len(provider.sent) != 0 {
		t.Fatalf("expected no delivery for standard event, got %d", len(provider.sent))
	}

	stored, _ := repo.List()
	if len(stored) != 1 {
		t.Fatalf("expected 1 stored event, got %d", len(stored))
	}
	if stored[0].AttentionClass != event.AttentionStandard {
		t.Fatalf("expected standard attention class")
	}
}

func TestPipelineDefaultsToStandardWhenNoRuleMatches(t *testing.T) {
	p, _, provider := newTestPipeline([]policy.Rule{{
		Name: "jira-only",
		When: policy.RuleCondition{Source: "jira"},
		Then: policy.RuleOutcome{AttentionClass: "critical", Reason: "jira alert"},
	}})

	result, err := p.Handle(event.Event{ID: "evt-3", Source: "github", EventType: "star"})
	if err != nil {
		t.Fatalf("handle event: %v", err)
	}

	if result.AttentionClass != event.AttentionStandard {
		t.Fatalf("expected standard default, got %s", result.AttentionClass)
	}
	if result.DecisionState.Status != event.DecisionQueued {
		t.Fatalf("expected queued, got %s", result.DecisionState.Status)
	}
	if len(provider.sent) != 0 {
		t.Fatalf("expected no delivery")
	}
}
