package storage

import (
	"context"
	"errors"

	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/mappers"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type taskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) ports.TaskRepo {
	return &taskRepo{
		db: db,
	}
}

func (r *taskRepo) GetListByBoardID(ctx context.Context, boardID uint, limit uint, offset uint) ([]domains.Task, uint, error) {
	var taskEntities []entities.Task

	query := r.db.WithContext(ctx).
		Model(&entities.Task{}).
		Where("board_id = ?", boardID).
		Preload("Board").
		Preload("Column").
		Preload("Assignee").
		Preload("Creator")

	//calculate total entities
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	//apply offset
	if offset > 0 {
		query = query.Offset(int(offset))
	}

	//apply limit
	if limit > 0 {
		query = query.Limit(int(limit))
	}

	//fetch entities
	if err := query.Find(&taskEntities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, fiber.NewError(fiber.StatusNotFound, "There is no task in the board!")
		}
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	taskModels := mappers.TaskEntitiesToDomain(taskEntities)

	return taskModels, uint(total), nil
}

func (r *taskRepo) Create(ctx context.Context, task *domains.Task) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newTask := mappers.DomainToTaskEntity(task)
		if err := tx.Create(&newTask).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		task.ID = newTask.ID
		return nil
	})

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (r *taskRepo) GetByID(ctx context.Context, id uint) (*domains.Task, error) {
	var task entities.Task
	err := r.db.WithContext(ctx).Model(&entities.Task{}).
		Where("id = ?", id).
		Preload("Board").
		Preload("Creator").
		Preload("Column").
		Preload("Assignee").
		Preload("Parent").
		First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Task not found!")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return mappers.TaskEntityToDomain(&task), nil
}

func (r *taskRepo) Update(ctx context.Context, task *domains.Task) error {
	var existingTask *entities.Task
	err := r.db.WithContext(ctx).Model(&entities.Task{}).Where("id = ?", task.ID).First(&existingTask).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Task not found!")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	existingTask.Name = task.Name
	existingTask.ParentID = task.ParentID
	existingTask.AssigneeID = task.AssigneeID
	existingTask.ColumnID = task.ColumnID
	existingTask.OrderPosition = task.OrderPosition
	existingTask.Name = task.Name
	existingTask.Description = task.Description
	existingTask.StartDateTime = task.StartDateTime
	existingTask.EndDateTime = task.EndDateTime
	existingTask.StoryPoint = task.StoryPoint

	if err := r.db.WithContext(ctx).Save(&existingTask).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *taskRepo) Delete(ctx context.Context, id uint) error {

	var existingTask *entities.Task
	err := r.db.WithContext(ctx).Model(&entities.Task{}).Where("id = ?", id).First(&existingTask).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Task not found!")
		}
		return err
	}

	if err := r.db.WithContext(ctx).Model(&entities.Task{}).Delete(&existingTask).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *taskRepo) GetTaskChildren(ctx context.Context, taskID uint) ([]domains.TaskChild, error) {
	var childEntities []entities.TaskChild
	query := `
        WITH RECURSIVE sub_tasks AS (SELECT *
									 FROM tasks
									 WHERE parent_id = ?
									 UNION ALL
									 SELECT t.*
									 FROM tasks t
											  INNER JOIN sub_tasks st ON st.id = t.parent_id)
		SELECT st2.*, columns.name AS "column_name", columns.is_final
		FROM sub_tasks st2
				 INNER JOIN columns on st2.column_id = columns.id
    `
	if err := r.db.WithContext(ctx).Raw(query, taskID).Scan(&childEntities).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	taskChildren := mappers.TaskChildEntitiesToDomain(childEntities)

	return taskChildren, nil
}

func (r *taskRepo) GetTaskDependencies(ctx context.Context, taskID uint) ([]domains.TaskDependency, error) {
	var dependencies []entities.TaskDependency
	if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Find(&dependencies).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return mappers.TaskDependencyEntitiesToDomains(dependencies), nil
}

func (r *taskRepo) AddTaskDependency(ctx context.Context, taskID, dependentTaskID uint) error {
	exists, err := r.DependencyExists(ctx, taskID, dependentTaskID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if exists {
		return fiber.NewError(fiber.StatusBadRequest, "Dependency already exists!")
	}
	dependency := entities.TaskDependency{TaskID: taskID, DependentTaskID: dependentTaskID}
	if err := r.db.WithContext(ctx).Create(&dependency).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *taskRepo) RemoveTaskDependency(ctx context.Context, taskID, dependentTaskID uint) error {
	exists, err := r.DependencyExists(ctx, taskID, dependentTaskID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if !exists {
		return fiber.NewError(fiber.StatusBadRequest, "Dependency already exists!")
	}
	if err := r.db.WithContext(ctx).Where("task_id = ? AND dependent_task_id = ?", taskID, dependentTaskID).Delete(&entities.TaskDependency{}).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *taskRepo) DependencyExists(ctx context.Context, taskID, dependentTaskID uint) (bool, error) {
	var dependencies []*entities.TaskDependency
	if err := r.db.WithContext(ctx).Where("task_id = ? AND dependent_task_id = ?", taskID, dependentTaskID).Find(&dependencies).Error; err != nil {
		return false, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return dependencies != nil, nil
}

func (r *taskRepo) GetAllTaskDependencies(ctx context.Context) ([]domains.TaskDependency, error) {
	var dependencies []entities.TaskDependency
	result := r.db.WithContext(ctx).Find(&dependencies).Error
	if result != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, result.Error())
	}
	return mappers.TaskDependencyEntitiesToDomains(dependencies), nil
}
