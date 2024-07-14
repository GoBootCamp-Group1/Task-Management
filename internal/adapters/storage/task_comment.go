package storage

import (
	"context"
	"errors"

	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/mappers"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type taskCommentRepo struct {
	db *gorm.DB
}

func NewTaskCommentRepo(db *gorm.DB) ports.TaskCommentRepo {
	return &taskCommentRepo{
		db: db,
	}
}

func (r *taskCommentRepo) Create(ctx context.Context, comment *domains.TaskComment) error {
	//generate UUID
	comment.ID = uuid.New()

	newComment := mappers.DomainToCommentEntity(comment)
	if err := r.db.Debug().WithContext(ctx).Create(&newComment).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	comment.ID = newComment.ID

	return nil
}

func (r *taskCommentRepo) GetByID(ctx context.Context, id string) (*domains.TaskComment, error) {
	var comment entities.TaskComment
	err := r.db.WithContext(ctx).Model(&entities.TaskComment{}).
		Where("id = ?", id).
		Preload("User").
		First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, "Comment not found!")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return mappers.TaskCommentEntityToDomain(&comment), nil
}

func (r *taskCommentRepo) Update(ctx context.Context, comment *domains.TaskComment) error {
	//TODO implement me
	panic("implement me")
}

func (r *taskCommentRepo) Delete(ctx context.Context, id string) error {
	var existingTaskComment *entities.TaskComment
	err := r.db.WithContext(ctx).Model(&entities.TaskComment{}).Where("id = ?", id).First(&existingTaskComment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Comment not found!")
		}
		return err
	}

	if err := r.db.WithContext(ctx).Model(&entities.TaskComment{}).Delete(&existingTaskComment).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (r *taskCommentRepo) GetListByTaskID(ctx context.Context, taskID uint, limit uint, offset uint) ([]domains.TaskComment, uint, error) {
	var commentEntities []entities.TaskComment

	query := r.db.WithContext(ctx).
		Model(&entities.TaskComment{}).
		Where("task_id = ?", taskID).
		Preload("User")

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
	if err := query.Find(&commentEntities).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, fiber.NewError(fiber.StatusNotFound, "There is no comment!")
		}
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	commentModels := mappers.TaskCommentEntitiesToDomain(commentEntities)

	return commentModels, uint(total), nil
}
