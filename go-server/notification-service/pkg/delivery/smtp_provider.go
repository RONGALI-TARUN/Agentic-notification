package delivery

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/example/notification-service/pkg/config"
	"github.com/example/notification-service/pkg/notification"
)

type SMTPProvider struct {
	cfg *config.Config
}

func NewSMTPProvider(cfg *config.Config) *SMTPProvider {
	return &SMTPProvider{cfg: cfg}
}

func (p *SMTPProvider) Send(req notification.NotificationRequest) error {
	if p.cfg == nil || p.cfg.SMTPHost == "" || p.cfg.SMTPSender == "" {
		return fmt.Errorf("smtp provider is not configured")
	}
	if len(req.EmailRecipients) == 0 {
		return nil
	}

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", strings.Join(req.EmailRecipients, ", "), req.Title, req.Message))
	auth := smtp.PlainAuth("", p.cfg.SMTPUsername, p.cfg.SMTPPassword, p.cfg.SMTPHost)
	return smtp.SendMail(fmt.Sprintf("%s:%d", p.cfg.SMTPHost, p.cfg.SMTPPort), auth, p.cfg.SMTPSender, req.EmailRecipients, msg)
}

func (p *SMTPProvider) Name() string {
	return "smtp"
}
