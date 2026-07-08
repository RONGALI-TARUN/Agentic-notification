package policy

import "github.com/example/notification-service/pkg/event"

type Engine struct {
	rules []Rule
}

func NewEngine(rules []Rule) *Engine {
	return &Engine{rules: rules}
}

func (e *Engine) Evaluate(evt event.Event) (PolicyDecision, bool) {
	for _, rule := range e.rules {
		if e.matches(rule.When, evt) {
			return PolicyDecision{
				AttentionClass: rule.Then.AttentionClass,
				Reason:         rule.Then.Reason,
				MatchedRule:    rule.Name,
			}, true
		}
	}
	return PolicyDecision{}, false
}

func (e *Engine) matches(condition RuleCondition, evt event.Event) bool {
	if condition.Source != "" && condition.Source != evt.Source {
		return false
	}
	if condition.EventType != "" && condition.EventType != evt.EventType {
		return false
	}
	if condition.Repository != "" && condition.Repository != evt.Repository {
		return false
	}
	return true
}
