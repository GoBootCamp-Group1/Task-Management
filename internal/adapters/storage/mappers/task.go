package mappers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"gorm.io/gorm"
)

func DomainToTaskEntity(task *domains.Task) *entities.Task {
	return &entities.Task{
		Model:     gorm.Model{ID: task.ID},
		CreatedBy: task.CreatedBy,
		Name:      task.Name,
	}
}

func TaskEntityToDomain(entity *entities.Task) *domains.Task {
	return &domains.Task{
		ID:        entity.ID,
		CreatedBy: entity.CreatedBy,
		Name:      entity.Name,
	}
}
