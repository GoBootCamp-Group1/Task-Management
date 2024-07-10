package notifier

import (
	"context"
	"errors"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/notification"
	notification2 "github.com/GoBootCamp-Group1/Task-Management/pkg/notification/notification"
)

type Adapter struct {
	ports.Notifier
	notifier *notification.Notifier
}

func NewNotifierAdapter(notifier *notification.Notifier) *Adapter {
	return &Adapter{notifier: notifier}
}

func (a *Adapter) SendInAppNotification(ctx context.Context, userID uint, input ports.NotificationInput) error {
	if a.notifier.InApp == nil {
		return errors.New("in-app notifier is not configured")
	}

	notificationInput := &notification2.InAppInput{
		Type:    input.Type,
		Message: input.Message,
	}

	return a.notifier.InApp.Send(ctx, userID, notificationInput)
}

func (a *Adapter) SendEmailNotification(to string, subject string, body string) error {
	if a.notifier.Email == nil {
		return errors.New("email notifier is not configured")
	}
	return a.notifier.Email.Send(to, subject, body)
}
