package services

import (
	"context"
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
)

type TaskService struct {
	repo            ports.TaskRepo
	taskCommentRepo ports.TaskCommentRepo
	notifier        ports.Notifier
	boardService    *BoardService
	columnService   *ColumnService
}

func NewTaskService(
	repo ports.TaskRepo,
	notifier ports.Notifier,
	boardService *BoardService,
	columnService *ColumnService,
	taskCommentRepo ports.TaskCommentRepo,
) *TaskService {
	return &TaskService{
		repo:            repo,
		notifier:        notifier,
		boardService:    boardService,
		columnService:   columnService,
		taskCommentRepo: taskCommentRepo,
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

	//Send notification to Assignee if exists
	if taskWithRelations.Assignee != nil {
		//Send notification
		input := ports.NotificationInput{
			Type:    ports.NewTaskAssignedNotification,
			Message: fmt.Sprintf("hey, %s. new task assigned.", taskWithRelations.Assignee.Name),
		}
		err := s.notifier.SendInAppNotification(ctx, taskWithRelations.Assignee.ID, input)
		if err != nil {
			return nil, err
		}
	}

	return taskWithRelations, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, userID uint, boardID uint, id uint) (*domains.Task, error) {
	task, errFetch := s.repo.GetByID(ctx, id)
	if errFetch != nil {
		return nil, fmt.Errorf("repository: can not fetch task #%d: %w", id, errFetch)
	}

	board, errFetchBoard := s.boardService.GetBoardByID(ctx, boardID)
	if errFetchBoard != nil {
		return nil, fmt.Errorf("board service: can not fetch board #%d: %w", boardID, errFetchBoard)
	}
	//check permission
	if board.IsPrivate {
		//check permissions -> only board members can see task
		hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Viewer, userID, boardID)
		if !hasAccess {
			return nil, fmt.Errorf("access denied")
		}
	}

	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, userID uint, boardID uint, task *domains.Task) (*domains.Task, error) {
	//check permissions
	hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Maintainer, userID, boardID)
	if !hasAccess {
		return nil, fmt.Errorf("access denied")
	}

	errUpdate := s.repo.Update(ctx, task)
	if errUpdate != nil {
		return nil, fmt.Errorf("repository: can not update task: %w", errUpdate)
	}

	taskWithRelations, errFetch := s.repo.GetByID(ctx, task.ID)
	if errFetch != nil {
		return nil, fmt.Errorf("repository: can not fetch task #%d %w", task.ID, errFetch)
	}

	return taskWithRelations, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, userID uint, id uint) error {
	//load task
	task, errFetch := s.repo.GetByID(ctx, id)
	if errFetch != nil {
		return fmt.Errorf("repository: can not fetch task #%d: %w", id, errFetch)
	}

	//check permissions
	hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Maintainer, userID, task.BoardID)
	if !hasAccess {
		return fmt.Errorf("access denied")
	}
	errDelete := s.repo.Delete(ctx, id)
	if errDelete != nil {
		return fmt.Errorf("repository: can not delete task %w", errDelete)
	}
	return nil
}

func (s *TaskService) ChangeTaskColumn(ctx context.Context, userID uint, task *domains.Task, newColumnID uint) (*domains.Task, error) {
	//check permissions
	hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Editor, userID, task.BoardID)
	if !hasAccess {
		return nil, fmt.Errorf("access denied")
	}

	//fetch task info
	t, errFetch := s.repo.GetByID(ctx, task.ID)
	if errFetch != nil {
		return nil, fmt.Errorf("repository: can not fetch task #%d %w", task.ID, errFetch)
	}

	newColumn, errFetchColumn := s.columnService.GetColumnById(ctx, newColumnID)
	if errFetchColumn != nil {
		return nil, fmt.Errorf("service: can not fetch column info #%d %w", newColumnID, errFetchColumn)
	}

	//Check for children tasks
	if newColumn.IsFinal {
		childrenTasks, errFetchChildrenTasks := s.repo.GetTaskChildren(ctx, task.ID)
		if errFetchChildrenTasks != nil {
			return nil, fmt.Errorf("repository: can not get children of task: %w", errFetchChildrenTasks)
		}

		allChildrenFinished := true
		for _, child := range childrenTasks {
			if !child.ColumnIsFinal {
				allChildrenFinished = false
				break
			}
		}

		if !allChildrenFinished {
			return nil, fmt.Errorf("service: children tasks are not finished yet")
		}
	}

	//change column and update
	t.ColumnID = newColumnID
	errUpdate := s.repo.Update(ctx, t)
	if errUpdate != nil {
		return nil, fmt.Errorf("repository: can not update task: %w", errUpdate)
	}

	taskWithRelations, errFetch := s.repo.GetByID(ctx, task.ID)
	if errFetch != nil {
		return nil, fmt.Errorf("repository: can not fetch task #%d %w", task.ID, errFetch)
	}

	return taskWithRelations, nil
}

func (s *TaskService) GetTaskChildren(ctx context.Context, userID uint, boardID uint, taskID uint) ([]domains.TaskChild, error) {
	childrenTasks, errFetchChildrenTasks := s.repo.GetTaskChildren(ctx, taskID)

	if errFetchChildrenTasks != nil {
		return nil, fmt.Errorf("repository: can not get children of task: %w", errFetchChildrenTasks)
	}

	return childrenTasks, nil
}

func (s *TaskService) CreateComment(ctx context.Context, userID uint, boardID uint, taskComment *domains.TaskComment) (*domains.TaskComment, error) {
	//check permissions
	hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Editor, userID, boardID)
	if !hasAccess {
		return nil, fmt.Errorf("access denied")
	}

	//create task comment
	errCreate := s.taskCommentRepo.Create(ctx, taskComment)
	if errCreate != nil {
		return nil, fmt.Errorf("repository: can not create comment: %w", errCreate)
	}

	//load comment
	comment, errFetch := s.taskCommentRepo.GetByID(ctx, taskComment.ID.String())
	if errFetch != nil {
		return nil, fmt.Errorf("repository: can not fetch comment #%d: %w", taskComment.ID, errFetch)
	}

	return comment, nil
}

func (s *TaskService) GetTaskComments(ctx context.Context, userID uint, boardID uint, taskID uint, pageNumber uint, pageSize uint) ([]domains.TaskComment, uint, error) {
	//check permissions
	hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Viewer, userID, boardID)
	if !hasAccess {
		return nil, 0, fmt.Errorf("access denied")
	}

	//pagination calculate
	limit := pageSize
	offset := (pageNumber - 1) * pageSize

	//create task comment
	comments, total, errFetch := s.taskCommentRepo.GetListByTaskID(ctx, taskID, limit, offset)
	if errFetch != nil {
		return nil, 0, fmt.Errorf("repository: can not fetch comments: %w", errFetch)
	}

	return comments, total, nil
}
