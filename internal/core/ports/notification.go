package ports

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
)

type NotificationRepo interface {
	GetByID(ctx context.Context, id string) (*domains.Notification, error)
	Read(ctx context.Context, notification *domains.Notification) (*domains.Notification, error)
	UnRead(ctx context.Context, notification *domains.Notification) (*domains.Notification, error)
	Delete(ctx context.Context, notification *domains.Notification) error
	GetList(ctx context.Context, userID uint, limit uint, offset uint) ([]domains.Notification, uint, error)
	GetUnreadList(ctx context.Context, userID uint, limit uint, offset uint) ([]domains.Notification, uint, error)
}
