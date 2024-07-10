package mappers

import (
	"database/sql"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/fp"
	"gorm.io/gorm"
	"time"
)

func NotificationEntityToDomain(entity *entities.Notification) *domains.Notification {
	var deletedAt, readAt *time.Time

	if entity.DeletedAt.Valid {
		deletedAt = &entity.DeletedAt.Time
	}

	if entity.ReadAt.Valid {
		readAt = &entity.ReadAt.Time
	}

	return &domains.Notification{
		ID:        entity.ID,
		UserID:    entity.UserID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: deletedAt,
		ReadAt:    readAt,
		Type:      entity.Type,
		Message:   entity.Message,

		User: UserEntityToDomain(&entity.User),
	}
}

func DomainToNotificationEntity(model *domains.Notification) *entities.Notification {
	var deletedAt, readAt sql.NullTime

	if model.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: *model.DeletedAt, Valid: true}
	} else {
		deletedAt = sql.NullTime{Valid: false}
	}

	if model.ReadAt != nil {
		readAt = sql.NullTime{Time: *model.ReadAt, Valid: true}
	} else {
		readAt = sql.NullTime{Valid: false}
	}
	return &entities.Notification{
		ID:        model.ID,
		UserID:    model.UserID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: gorm.DeletedAt(deletedAt),
		ReadAt:    readAt,
		Type:      model.Type,
		Message:   model.Message,
	}
}

func NotificationEntitiesToDomain(notificationEntities []entities.Notification) []domains.Notification {
	return fp.Map(notificationEntities, func(entity entities.Notification) domains.Notification {
		return *NotificationEntityToDomain(&entity)
	})
}

func NotificationDomainsToEntity(notificationDomains []domains.Notification) []entities.Notification {
	return fp.Map(notificationDomains, func(member domains.Notification) entities.Notification {
		return *DomainToNotificationEntity(&member)
	})
}
