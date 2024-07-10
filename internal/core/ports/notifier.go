package ports

import (
	"context"
)

type NotificationInput struct {
	Type string
	Data string
}

type Notifier interface {
	SendInAppNotification(ctx context.Context, userID uint, input NotificationInput) error
	SendEmailNotification(to string, subject string, body string) error
}
