package notification

type AgentInfo struct {
	Name  string `json:"name"`
	Owner string `json:"owner,omitempty"`
	Group string `json:"group,omitempty"`
}

type NotificationRequest struct {
	Agent            AgentInfo `json:"agent"`
	Title            string    `json:"title"`
	Message          string    `json:"message"`
	EmailRecipients  []string  `json:"emailRecipients,omitempty"`
	SlackWebhookURLs []string  `json:"slackWebhookUrls,omitempty"`
	TeamsWebhookURLs []string  `json:"teamsWebhookUrls,omitempty"`
}
