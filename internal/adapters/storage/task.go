package storage

import (
	"context"
	"errors"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/mappers"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
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

var (
	ErrTaskAlreadyExists = errors.New("task already exists")
)

func (r *taskRepo) Create(ctx context.Context, task *domains.Task) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newTask := mappers.DomainToTaskEntity(task)
		if err := tx.Create(&newTask).Error; err != nil {
			return err
		}
		task.ID = newTask.ID
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *taskRepo) GetByID(ctx context.Context, id uint) (*domains.Task, error) {
	var task entities.Task
	err := r.db.WithContext(ctx).Model(&entities.Task{}).
		Where("id = ?", id).
		Preload("Board").
		Preload("Creator").
		//TODO:additional relations
		First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mappers.TaskEntityToDomain(&task), nil
}

func (r *taskRepo) Update(ctx context.Context, task *domains.Task) error {
	var existingTask *entities.Task
	err := r.db.WithContext(ctx).Model(&entities.Task{}).Where("id = ?", task.ID).First(&existingTask).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
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
	existingTask.Additional = task.Additional

	return r.db.WithContext(ctx).Save(&existingTask).Error
}

func (r *taskRepo) Delete(ctx context.Context, id uint) error {

	var existingTask *entities.Task
	err := r.db.WithContext(ctx).Model(&entities.Task{}).Where("id = ?", id).First(&existingTask).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return r.db.WithContext(ctx).Model(&entities.Task{}).Delete(&existingTask).Error
}
