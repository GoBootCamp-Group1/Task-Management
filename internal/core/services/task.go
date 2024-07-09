package services

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
)

type TaskService struct {
	repo ports.TaskRepo
}

func NewTaskService(repo ports.TaskRepo) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, task *domains.Task) error {
	return s.repo.Create(ctx, task)
}

func (s *TaskService) GetTaskByID(ctx context.Context, id uint) (*domains.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) UpdateTask(ctx context.Context, task *domains.Task) error {
	return s.repo.Update(ctx, task)
}

func (s *TaskService) DeleteTask(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
