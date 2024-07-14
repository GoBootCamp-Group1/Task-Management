package services

import (
	"context"

	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

type NotificationService struct {
	repo ports.NotificationRepo
}

func NewNotificationService(repo ports.NotificationRepo) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

func (s *NotificationService) GetAllNotificationsList(ctx context.Context, userID uint, pageNumber uint, pageSize uint) ([]domains.Notification, uint, error) {
	limit := pageSize
	offset := (pageNumber - 1) * pageSize

	//fetch notifications
	notifications, total, errFetch := s.repo.GetList(ctx, userID, limit, offset)
	if errFetch != nil {
		return nil, 0, errFetch
	}

	return notifications, total, nil
}

func (s *NotificationService) GetUnReadNotificationsList(ctx context.Context, userID uint, pageNumber uint, pageSize uint) ([]domains.Notification, uint, error) {
	limit := pageSize
	offset := (pageNumber - 1) * pageSize

	//fetch notifications
	notifications, total, errFetch := s.repo.GetUnreadList(ctx, userID, limit, offset)
	if errFetch != nil {
		return nil, 0, errFetch
	}

	return notifications, total, nil
}

func (s *NotificationService) GetNotificationByID(ctx context.Context, id string, userID uint) (*domains.Notification, error) {
	notification, errFetch := s.repo.GetByID(ctx, id)
	if errFetch != nil {
		return nil, errFetch
	}

	if notification.UserID != userID {
		return nil, &fiber.Error{Code: fiber.StatusForbidden, Message: "Access denied"}
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
		return nil, errRead
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
		return nil, errRead
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
		return errDelete
	}
	return nil
}
