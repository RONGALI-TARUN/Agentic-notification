package feedback

type EventFeedback struct {
	ID        string `json:"id"`
	EventID   string `json:"eventId"`
	UserID    string `json:"userId,omitempty"`
	Action    string `json:"action"`
	Reason    string `json:"reason,omitempty"`
	CreatedAt string `json:"createdAt"`
}
