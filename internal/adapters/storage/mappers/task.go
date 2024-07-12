package mappers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/fp"
	"gorm.io/gorm"
)

func DomainToTaskEntity(model *domains.Task) *entities.Task {
	return &entities.Task{
		Model:         gorm.Model{ID: model.ID},
		CreatedBy:     model.CreatedBy,
		BoardID:       model.BoardID,
		ParentID:      model.ParentID,
		AssigneeID:    model.AssigneeID,
		ColumnID:      model.ColumnID,
		OrderPosition: model.OrderPosition,
		Name:          model.Name,
		Description:   model.Description,
		StartDateTime: model.StartDateTime,
		EndDateTime:   model.EndDateTime,
		StoryPoint:    model.StoryPoint,
		Additional:    model.Additional,
	}
}

func TaskEntityToDomain(entity *entities.Task) *domains.Task {
	return &domains.Task{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		CreatedBy: entity.CreatedBy,

		BoardID:       entity.BoardID,
		ParentID:      entity.ParentID,
		AssigneeID:    entity.AssigneeID,
		ColumnID:      entity.ColumnID,
		OrderPosition: entity.OrderPosition,
		Name:          entity.Name,
		Description:   entity.Description,
		StartDateTime: entity.StartDateTime,
		EndDateTime:   entity.EndDateTime,
		StoryPoint:    entity.StoryPoint,
		Additional:    entity.Additional,

		Board:   BoardEntityToDomain(&entity.Board),
		Creator: UserEntityToDomain(&entity.Creator),
		Column:  ColumnEntityToDomain(&entity.Column),
	}
}

func TaskEntitiesToDomain(taskEntities []entities.Task) []domains.Task {
	return fp.Map(taskEntities, func(entity entities.Task) domains.Task {
		return *TaskEntityToDomain(&entity)
	})
}

func TaskDomainsToEntity(taskDomains []domains.Task) []entities.Task {
	return fp.Map(taskDomains, func(member domains.Task) entities.Task {
		return *DomainToTaskEntity(&member)
	})
}
