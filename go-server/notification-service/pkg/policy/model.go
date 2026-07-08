package policy

type RuleCondition struct {
	Source     string `json:"source,omitempty"`
	EventType  string `json:"eventType,omitempty"`
	Repository string `json:"repository,omitempty"`
}

type RuleOutcome struct {
	AttentionClass string `json:"attentionClass"`
	Reason         string `json:"reason"`
}

type Rule struct {
	Name string        `json:"name"`
	When RuleCondition `json:"when"`
	Then RuleOutcome   `json:"then"`
}

type PolicyDecision struct {
	AttentionClass string
	Reason         string
	MatchedRule    string
}
