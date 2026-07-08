package delivery

import (
	"testing"

	"github.com/example/notification-service/pkg/config"
	"github.com/example/notification-service/pkg/notification"
)

func TestSMTPProviderRejectsUnconfiguredState(t *testing.T) {
	provider := NewSMTPProvider(&config.Config{})
	if err := provider.Send(notification.NotificationRequest{EmailRecipients: []string{"test@example.com"}}); err == nil {
		t.Fatal("expected unconfigured smtp provider to fail")
	}
}
