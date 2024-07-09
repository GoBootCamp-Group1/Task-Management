package storage

import (
	"context"
	"errors"
	"fmt"
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
	//var newTask entities.Task

	fmt.Println(task)

	return r.db.Transaction(func(tx *gorm.DB) error {
		fmt.Println("reached")
		//entity := mappers.DomainToTaskEntity(task)
		//
		//if err := tx.WithContext(ctx).Create(&entity).Error; err != nil {
		//	return err
		//}
		//task.ID = entity.ID // set the ID to the domain object after creation
		return nil
	})
}

func (r *taskRepo) GetByID(ctx context.Context, id uint) (*domains.Task, error) {
	var b entities.Task
	err := r.db.WithContext(ctx).Model(&entities.Task{}).Where("id = ?", id).First(&b).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mappers.TaskEntityToDomain(&b), nil
}

func (r *taskRepo) Update(ctx context.Context, task *domains.Task) error {
	var existingTask *entities.Task
	if err := r.db.WithContext(ctx).Model(&entities.Task{}).Where("id = ?", task.ID).First(&existingTask).Error; err != nil {
		return err
	}

	existingTask.Name = task.Name

	return r.db.WithContext(ctx).Save(&existingTask).Error
}

func (r *taskRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Task{}, id).Error
}
