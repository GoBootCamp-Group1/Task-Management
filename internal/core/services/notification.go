package services

import (
	"context"
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
)

type NotificationService struct {
	repo ports.NotificationRepo
}

func NewNotificationService(repo ports.NotificationRepo) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

func (s *NotificationService) GetAllNotificationsList(ctx context.Context, userID uint, pageNumber uint, pageSize uint) ([]*domains.Notification, uint, error) {
	limit := pageSize
	offset := (pageNumber - 1) * pageSize

	//fetch notifications
	notifications, total, errFetch := s.repo.GetList(ctx, userID, limit, offset)
	if errFetch != nil {
		return nil, 0, fmt.Errorf("repository: can not fetch notifications: %w", errFetch)
	}

	return notifications, total, nil
}

func (s *NotificationService) GetUnReadNotificationsList(ctx context.Context, userID uint, pageNumber uint, pageSize uint) ([]*domains.Notification, uint, error) {
	limit := pageSize
	offset := (pageNumber - 1) * pageSize

	//fetch notifications
	notifications, total, errFetch := s.repo.GetUnreadList(ctx, userID, limit, offset)
	if errFetch != nil {
		return nil, 0, fmt.Errorf("repository: can not fetch notifications: %w", errFetch)
	}

	return notifications, total, nil
}

func (s *NotificationService) GetNotificationByID(ctx context.Context, id string, userID uint) (*domains.Notification, error) {
	notification, errFetch := s.repo.GetByID(ctx, id)
	if errFetch != nil {
		return nil, fmt.Errorf("repository: can not fetch notification #%s %w", id, errFetch)
	}

	if notification.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	return notification, nil
}

func (s *NotificationService) ReadNotification(ctx context.Context, id string, userID uint) (*domains.Notification, error) {
	notification, errFetch := s.GetNotificationByID(ctx, id, userID)
	if errFetch != nil {
		return nil, errFetch
	}

	readNotification, errRead := s.repo.Read(ctx, notification)
	if errRead != nil {
		return nil, fmt.Errorf("repository: can not read notification: %w", errRead)
	}

	return readNotification, nil
}

func (s *NotificationService) UnReadNotification(ctx context.Context, id string, userID uint) (*domains.Notification, error) {
	notification, errFetch := s.GetNotificationByID(ctx, id, userID)
	if errFetch != nil {
		return nil, errFetch
	}

	unreadNotification, errRead := s.repo.UnRead(ctx, notification)
	if errRead != nil {
		return nil, fmt.Errorf("repository: can not unread notification: %w", errRead)
	}

	return unreadNotification, nil
}

func (s *NotificationService) DeleteNotification(ctx context.Context, id string, userID uint) error {
	notification, errFetch := s.GetNotificationByID(ctx, id, userID)
	if errFetch != nil {
		return errFetch
	}

	errDelete := s.repo.Delete(ctx, notification)
	if errDelete != nil {
		return fmt.Errorf("repository: can not delete notification %w", errDelete)
	}
	return nil
}
