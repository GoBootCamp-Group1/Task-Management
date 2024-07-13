package mappers

import (
	"database/sql"
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

	var assignee *domains.User
	if entity.AssigneeID != nil {
		assignee = UserEntityToDomain(entity.Assignee)
	}

	//var parent *domains.Task
	//if entity.ParentID != nil {
	//	parent = TaskEntityToDomain(entity.Parent)
	//}

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
		//Parent:   parent,
		Assignee: assignee,
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

func TaskChildEntityToDomain(entity *entities.TaskChild) *domains.TaskChild {
	return &domains.TaskChild{
		ID:            entity.ID,
		ColumnID:      entity.ColumnID,
		OrderPosition: entity.OrderPosition,
		Name:          entity.Name,
		Description:   entity.Description,
		ColumnName:    entity.ColumnName,
		ColumnIsFinal: entity.ColumnIsFinal,
	}
}

func TaskChildEntitiesToDomain(taskEntities []entities.TaskChild) []domains.TaskChild {
	return fp.Map(taskEntities, func(entity entities.TaskChild) domains.TaskChild {
		return *TaskChildEntityToDomain(&entity)
	})
}

func DomainToCommentEntity(model *domains.TaskComment) *entities.TaskComment {
	var deletedAt sql.NullTime

	if model.DeletedAt != nil {
		deletedAt = sql.NullTime{Time: *model.DeletedAt, Valid: true}
	} else {
		deletedAt = sql.NullTime{Valid: false}
	}

	return &entities.TaskComment{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: gorm.DeletedAt(deletedAt),
		UserID:    model.UserID,
		TaskID:    model.TaskID,
		Comment:   model.Comment,
	}
}

func TaskCommentEntityToDomain(entity *entities.TaskComment) *domains.TaskComment {
	return &domains.TaskComment{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		Comment:   entity.Comment,
		User:      UserEntityToDomain(&entity.User),
	}
}
