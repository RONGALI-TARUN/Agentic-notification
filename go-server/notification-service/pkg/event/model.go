package event

type Event struct {
	ID         string                 `json:"id"`
	Source     string                 `json:"source"`
	EventType  string                 `json:"eventType"`
	Repository string                 `json:"repository,omitempty"`
	Payload    map[string]interface{} `json:"payload"`
}
