package delivery

import (
	"fmt"

	"github.com/example/notification-service/pkg/notification"
)

type Service struct {
	provider Provider
}

func NewService(provider Provider) *Service {
	return &Service{provider: provider}
}

func (s *Service) Send(req notification.NotificationRequest) error {
	if s.provider == nil {
		return fmt.Errorf("no delivery provider configured")
	}
	return s.provider.Send(req)
}

func (s *Service) ProviderName() string {
	if s.provider == nil {
		return "none"
	}
	return s.provider.Name()
}
