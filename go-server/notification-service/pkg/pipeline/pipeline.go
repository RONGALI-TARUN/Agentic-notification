package pipeline

import (
	"fmt"

	"github.com/example/notification-service/pkg/delivery"
	"github.com/example/notification-service/pkg/event"
	"github.com/example/notification-service/pkg/notification"
	"github.com/example/notification-service/pkg/policy"
)

type Config struct {
	DefaultEmailRecipients []string
	SlackWebhookURLs       []string
}

type Pipeline struct {
	policyEngine *policy.Engine
	deliverySvc  *delivery.Service
	eventRepo    event.Repository
	cfg          Config
}

func New(
	engine *policy.Engine,
	deliverySvc *delivery.Service,
	eventRepo event.Repository,
	cfg Config,
) *Pipeline {
	return &Pipeline{
		policyEngine: engine,
		deliverySvc:  deliverySvc,
		eventRepo:    eventRepo,
		cfg:          cfg,
	}
}

func (p *Pipeline) Handle(evt event.Event) (event.Event, error) {
	evt = p.applyPolicy(evt)

	switch evt.AttentionClass {
	case event.AttentionCritical:
		if err := p.deliverCritical(evt); err != nil {
			return evt, fmt.Errorf("deliver critical event: %w", err)
		}
		evt.DecisionState.Status = event.DecisionSent
	default:
		if evt.DecisionState.Status == "" {
			evt.DecisionState.Status = event.DecisionQueued
		}
	}

	if err := p.eventRepo.Save(evt); err != nil {
		return evt, fmt.Errorf("persist event: %w", err)
	}
	return evt, nil
}

func (p *Pipeline) applyPolicy(evt event.Event) event.Event {
	decision, matched := p.policyEngine.Evaluate(evt)
	if matched {
		evt.AttentionClass = decision.AttentionClass
		evt.DecisionState = &event.DecisionState{
			Reason:      decision.Reason,
			MatchedRule: decision.MatchedRule,
		}
	} else {
		evt.AttentionClass = event.AttentionStandard
		evt.DecisionState = &event.DecisionState{
			Reason: "no matching rule, default standard",
		}
	}
	return evt
}

func (p *Pipeline) deliverCritical(evt event.Event) error {
	title := fmt.Sprintf("[CRITICAL] %s/%s", evt.Source, evt.EventType)
	message := fmt.Sprintf(
		"Critical alert\nSource: %s\nEvent: %s\nRepository: %s\nReason: %s",
		evt.Source, evt.EventType, evt.Repository, evt.DecisionState.Reason,
	)

	return p.deliverySvc.Send(notification.NotificationRequest{
		Agent: notification.AgentInfo{
			Name:  "attention-orchestration",
			Group: evt.Source,
		},
		Title:            title,
		Message:          message,
		EmailRecipients:  p.cfg.DefaultEmailRecipients,
		SlackWebhookURLs: p.cfg.SlackWebhookURLs,
	})
}
