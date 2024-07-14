package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/mappers"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type notificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) ports.NotificationRepo {
	return &notificationRepo{
		db: db,
	}
}

func (r *notificationRepo) GetByID(ctx context.Context, id string) (*domains.Notification, error) {
	var n entities.Notification
	err := r.db.WithContext(ctx).Model(&entities.Notification{}).Where("id = ?", id).First(&n).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Notification not found!")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return mappers.NotificationEntityToDomain(&n), nil
}

func (r *notificationRepo) Read(ctx context.Context, notification *domains.Notification) (*domains.Notification, error) {
	n := mappers.DomainToNotificationEntity(notification)
	n.ReadAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	err := r.db.WithContext(ctx).Save(&n).Error

	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return mappers.NotificationEntityToDomain(n), nil
}

func (r *notificationRepo) UnRead(ctx context.Context, notification *domains.Notification) (*domains.Notification, error) {
	n := mappers.DomainToNotificationEntity(notification)
	n.ReadAt = sql.NullTime{Valid: false}
	err := r.db.WithContext(ctx).Save(&n).Error

	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return mappers.NotificationEntityToDomain(n), nil
}

func (r *notificationRepo) Delete(ctx context.Context, notification *domains.Notification) error {
	n := mappers.DomainToNotificationEntity(notification)
	if err := r.db.WithContext(ctx).Delete(&n).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *notificationRepo) GetList(ctx context.Context, userID uint, limit uint, offset uint) ([]domains.Notification, uint, error) {
	var notificationEntities []entities.Notification

	query := r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("user_id = ?", userID).
		Preload("User")

	//calculate total entities
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	//apply offset
	if offset > 0 {
		query = query.Offset(int(offset))
	}

	//apply limit
	if limit > 0 {
		query = query.Limit(int(limit))
	}

	//fetch entities
	if err := query.Find(&notificationEntities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, nil
		}
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return mappers.NotificationEntitiesToDomain(notificationEntities), uint(total), nil
}

func (r *notificationRepo) GetUnreadList(ctx context.Context, userID uint, limit uint, offset uint) ([]domains.Notification, uint, error) {
	var notificationEntities []entities.Notification

	query := r.db.WithContext(ctx).
		Model(&entities.Notification{}).
		Where("read_at IS NULL").
		Where("user_id = ?", userID).
		Preload("User")

	//calculate total entities
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	//apply offset
	if offset > 0 {
		query = query.Offset(int(offset))
	}

	//apply limit
	if limit > 0 {
		query = query.Limit(int(limit))
	}

	//fetch entities
	if err := query.Find(&notificationEntities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, fiber.NewError(fiber.StatusNotFound, "There is no unread notification!")
		}
		return nil, 0, err
	}

	return mappers.NotificationEntitiesToDomain(notificationEntities), uint(total), nil
}
