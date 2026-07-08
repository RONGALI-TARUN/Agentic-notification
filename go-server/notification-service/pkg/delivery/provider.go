package delivery

import "github.com/example/notification-service/pkg/notification"

type Provider interface {
	Send(req notification.NotificationRequest) error
	Name() string
}

type StubProvider struct{}

func (s *StubProvider) Send(req notification.NotificationRequest) error {
	return nil
}

func (s *StubProvider) Name() string {
	return "stub"
}
