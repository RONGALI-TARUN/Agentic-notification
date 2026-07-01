package impl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/example/notification-service/pkg/config"
	"github.com/example/notification-service/pkg/notification"
)

type NotificationServiceImpl struct {
	cfg *config.Config
}

func NewNotificationService(cfg *config.Config) *NotificationServiceImpl {
	return &NotificationServiceImpl{cfg: cfg}
}

func (s *NotificationServiceImpl) Send(req notification.NotificationRequest) error {
	if req.Agent.Name == "" {
		return errors.New("agent name required")
	}
	if req.Title == "" {
		return errors.New("title required")
	}

	if len(req.EmailRecipients) > 0 {
		if err := s.sendEmail(req.EmailRecipients, req.Title, req.Message, req.Agent); err != nil {
			return err
		}
	}

	if len(req.SlackWebhookURLs) > 0 {
		for _, url := range req.SlackWebhookURLs {
			if err := sendWebhook(url, req.Message, req.Agent); err != nil {
				return err
			}
		}
	}

	if len(req.TeamsWebhookURLs) > 0 {
		for _, url := range req.TeamsWebhookURLs {
			if err := sendWebhook(url, req.Message, req.Agent); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *NotificationServiceImpl) sendEmail(to []string, title, message string, agent notification.AgentInfo) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)
	headers := make(map[string]string)
	headers["From"] = s.cfg.SMTPSender
	headers["To"] = strings.Join(to, ", ")
	headers["Subject"] = fmt.Sprintf("[n8n agent] %s", title)
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=UTF-8"

	var msg bytes.Buffer
	for k, v := range headers {
		fmt.Fprintf(&msg, "%s: %s\r\n", k, v)
	}
	msg.WriteString("\r\n")
	msg.WriteString(fmt.Sprintf("%s\n\nAgent: %s\nOwner: %s\nGroup: %s", message, agent.Name, agent.Owner, agent.Group))

	auth := smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)
	return smtp.SendMail(addr, auth, s.cfg.SMTPSender, to, msg.Bytes())
}

func sendWebhook(url, message string, agent notification.AgentInfo) error {
	payload, _ := json.Marshal(map[string]string{"text": fmt.Sprintf("%s\n\nAgent: %s\nGroup: %s", message, agent.Name, agent.Group)})
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook status %d", resp.StatusCode)
	}
	return nil
}
