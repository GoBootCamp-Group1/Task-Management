package mappers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/fp"
)

func DomainToTaskDependencyEntity(model domains.TaskDependency) *entities.TaskDependency {
	return &entities.TaskDependency{
		TaskID:          model.TaskID,
		DependentTaskID: model.DependentTaskID,
	}
}

func TaskDependencyEntityToDomain(entity entities.TaskDependency) *domains.TaskDependency {
	return &domains.TaskDependency{
		TaskID:          entity.TaskID,
		DependentTaskID: entity.DependentTaskID,
	}

}

func TaskDependencyEntitiesToDomains(taskEntities []entities.TaskDependency) []domains.TaskDependency {
	return fp.Map(taskEntities, func(entity entities.TaskDependency) domains.TaskDependency {
		return *TaskDependencyEntityToDomain(entity)
	})
}

func TaskDependencyDomainsToEntities(taskDomains []domains.TaskDependency) []entities.TaskDependency {
	return fp.Map(taskDomains, func(member domains.TaskDependency) entities.TaskDependency {
		return *DomainToTaskDependencyEntity(member)
	})
}
