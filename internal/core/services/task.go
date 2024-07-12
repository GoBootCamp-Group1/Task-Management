package services

import (
	"context"
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
)

type TaskService struct {
	repo         ports.TaskRepo
	notifier     ports.Notifier
	boardService *BoardService
}

func NewTaskService(repo ports.TaskRepo, notifier ports.Notifier, boardService *BoardService) *TaskService {
	return &TaskService{
		repo:         repo,
		notifier:     notifier,
		boardService: boardService,
	}
}

func (s *TaskService) GetTasksByBoardID(ctx context.Context, userID uint, boardID uint, pageNumber uint, pageSize uint) ([]domains.Task, uint, error) {
	//check permission
	board, errFetchBoard := s.boardService.GetBoardByID(ctx, boardID)
	if errFetchBoard != nil {
		return nil, 0, fmt.Errorf("board service: can not fetch board #%d: %w", boardID, errFetchBoard)
	}

	if board.IsPrivate {
		//check permissions -> only board members can see task
		hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Viewer, userID, boardID)
		if !hasAccess {
			return nil, 0, fmt.Errorf("access denied")
		}
	}

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
	//check permissions
	hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Maintainer, task.CreatedBy, task.BoardID)
	if !hasAccess {
		return nil, fmt.Errorf("access denied")
	}

	//create task
	errCreate := s.repo.Create(ctx, task)
	if errCreate != nil {
		return nil, fmt.Errorf("repository: can not create task: %w", errCreate)
	}

	//load task
	taskWithRelations, errFetch := s.repo.GetByID(ctx, task.ID)
	if errFetch != nil {
		return nil, fmt.Errorf("repository: can not fetch task #%d: %w", task.ID, errFetch)
	}

	return taskWithRelations, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, id uint) (*domains.Task, error) {
	task, errFetch := s.repo.GetByID(ctx, id)
	if errFetch != nil {
		return nil, fmt.Errorf("repository: can not fetch task #%d: %w", id, errFetch)
	}

	board, errFetchBoard := s.boardService.GetBoardByID(ctx, task.BoardID)
	if errFetchBoard != nil {
		return nil, fmt.Errorf("board service: can not fetch board #%d: %w", task.BoardID, errFetchBoard)
	}
	//check permission
	if board.IsPrivate {
		//check permissions -> only board members can see task
		hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Viewer, task.CreatedBy, task.BoardID)
		if !hasAccess {
			return nil, fmt.Errorf("access denied")
		}
	}

	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, task *domains.Task) (*domains.Task, error) {
	//check permissions
	hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Maintainer, task.CreatedBy, task.BoardID)
	if !hasAccess {
		return nil, fmt.Errorf("access denied")
	}

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
	//load task
	task, errFetch := s.repo.GetByID(ctx, id)
	if errFetch != nil {
		return fmt.Errorf("repository: can not fetch task #%d: %w", id, errFetch)
	}

	//check permissions
	hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Maintainer, task.CreatedBy, task.BoardID)
	if !hasAccess {
		return fmt.Errorf("access denied")
	}
	errDelete := s.repo.Delete(ctx, id)
	if errDelete != nil {
		return fmt.Errorf("repository: can not delete task %w", errDelete)
	}
	return nil
}
