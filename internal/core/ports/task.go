package ports

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
)

type TaskRepo interface {
	Create(ctx context.Context, task *domains.Task) error
	GetByID(ctx context.Context, id uint) (*domains.Task, error)
	Update(ctx context.Context, task *domains.Task) error
	Delete(ctx context.Context, id uint) error
	GetListByBoardID(ctx context.Context, boardID uint, limit uint, offset uint) ([]domains.Task, uint, error)
	GetTaskDependencies(ctx context.Context, taskID uint) ([]*entities.TaskDependency, error)
	AddTaskDependency(ctx context.Context, taskID, dependentTaskID uint) error
	RemoveTaskDependency(ctx context.Context, taskID, dependentTaskID uint) error
	DependencyExists(ctx context.Context, taskID, dependentTaskID uint) (bool, error)
	GetAllTaskDependencies(ctx context.Context) ([]entities.TaskDependency, error)
}
