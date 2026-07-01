package config

import (
	"os"
	"strconv"
)

type Config struct {
	Address         string
	SMTPHost        string
	SMTPPort        int
	SMTPUsername    string
	SMTPPassword    string
	SMTPSender      string
	SlackWebhookURL string
	TeamsWebhookURL string
}

func LoadConfig() (*Config, error) {
	port := 8080
	if v := os.Getenv("PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			port = p
		}
	}

	smtpPort := 587
	if v := os.Getenv("SMTP_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			smtpPort = p
		}
	}

	return &Config{
		Address:         ":" + strconv.Itoa(port),
		SMTPHost:        os.Getenv("SMTP_HOST"),
		SMTPPort:        smtpPort,
		SMTPUsername:    os.Getenv("SMTP_USERNAME"),
		SMTPPassword:    os.Getenv("SMTP_PASSWORD"),
		SMTPSender:      os.Getenv("SMTP_SENDER"),
		SlackWebhookURL: os.Getenv("SLACK_WEBHOOK_URL"),
		TeamsWebhookURL: os.Getenv("TEAMS_WEBHOOK_URL"),
	}, nil
}
