package services

import (
	"context"
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
)

type TaskService struct {
	repo ports.TaskRepo
}

func NewTaskService(repo ports.TaskRepo) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetTasksByBoardID(ctx context.Context, boardID uint, pageNumber uint, pageSize uint) ([]*domains.Task, uint, error) {
	//TODO: check for permissions
	//pagination calculate
	limit := pageSize
	offset := (pageNumber - 1) * pageSize

	//fetch tasks
	tasks, total, errFetch := s.repo.GetListByBoardID(ctx, boardID, limit, offset)
	if errFetch != nil {
		return nil, 0, fmt.Errorf("repository: can not fetch tasks: %w", errFetch)
	}

	return tasks, total, nil
}

func (s *TaskService) CreateTask(ctx context.Context, task *domains.Task) (*domains.Task, error) {
	//TODO: check for permissions
	//some business logic checks before creating task
	//check role
	//edit and move relate to assignee or maintainer and owner
	//get parent id

	errCreate := s.repo.Create(ctx, task)
	if errCreate != nil {
		return nil, fmt.Errorf("repository: can not create task: %w", errCreate)
	}

	taskWithRelations, errFetch := s.repo.GetByID(ctx, task.ID)
	if errFetch != nil {
		return nil, fmt.Errorf("repository: can not fetch task #%d: %w", task.ID, errFetch)
	}

	return taskWithRelations, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, id uint) (*domains.Task, error) {
	//TODO: check for permissions
	task, errFetch := s.repo.GetByID(ctx, id)
	if errFetch != nil {
		return nil, fmt.Errorf("repository: can not fetch task #%d: %w", id, errFetch)
	}
	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, task *domains.Task) (*domains.Task, error) {
	//TODO: check for permissions
	errUpdate := s.repo.Update(ctx, task)
	if errUpdate != nil {
		return nil, fmt.Errorf("repository: can not update task: %w", errUpdate)
	}

	taskWithRelations, errFetch := s.repo.GetByID(ctx, task.ID)
	if errFetch != nil {
		return nil, fmt.Errorf("repository: can not fetch task #%d %w", task.ID, errUpdate)
	}

	return taskWithRelations, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id uint) error {
	//TODO: check for permissions
	errDelete := s.repo.Delete(ctx, id)
	if errDelete != nil {
		return fmt.Errorf("repository: can not delete task %w", errDelete)
	}
	return nil
}
