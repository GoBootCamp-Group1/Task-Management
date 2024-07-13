package ports

import (
	"context"

	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
)

type TaskRepo interface {
	Create(ctx context.Context, task *domains.Task) error
	GetByID(ctx context.Context, id uint) (*domains.Task, error)
	Update(ctx context.Context, task *domains.Task) error
	Delete(ctx context.Context, id uint) error
	GetListByBoardID(ctx context.Context, boardID uint, limit uint, offset uint) ([]domains.Task, uint, error)
	GetTaskDependencies(ctx context.Context, taskID uint) ([]domains.TaskDependency, error)
	AddTaskDependency(ctx context.Context, taskID, dependentTaskID uint) error
	RemoveTaskDependency(ctx context.Context, taskID, dependentTaskID uint) error
	DependencyExists(ctx context.Context, taskID, dependentTaskID uint) (bool, error)
	GetAllTaskDependencies(ctx context.Context) ([]domains.TaskDependency, error)
	GetTaskChildren(ctx context.Context, taskID uint) ([]domains.TaskChild, error)
}

type TaskCommentRepo interface {
	Create(ctx context.Context, comment *domains.TaskComment) error
	GetByID(ctx context.Context, id string) (*domains.TaskComment, error)
	Update(ctx context.Context, comment *domains.TaskComment) error
	Delete(ctx context.Context, id string) error
	GetListByTaskID(ctx context.Context, taskID uint, limit uint, offset uint) ([]domains.TaskComment, uint, error)
}
