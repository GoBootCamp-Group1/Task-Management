package presenter

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"time"
)

type NotificationPresenter struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at"`
	Type      string     `json:"type"`
	Message   string     `json:"message"`
}

func NewNotificationPresenter(notification *domains.Notification) *NotificationPresenter {
	return &NotificationPresenter{
		ID:        notification.ID.String(),
		CreatedAt: notification.CreatedAt,
		ReadAt:    notification.ReadAt,
		Type:      notification.Type,
		Message:   notification.Message,
	}
}
