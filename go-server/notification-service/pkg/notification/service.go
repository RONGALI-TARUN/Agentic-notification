package notification

type NotificationService interface {
	Send(req NotificationRequest) error
}
