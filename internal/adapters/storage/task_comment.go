package storage

import (
	"context"
	"errors"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/mappers"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
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
		return err
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
			return nil, nil
		}
		return nil, err
	}
	return mappers.TaskCommentEntityToDomain(&comment), nil
}

func (r *taskCommentRepo) Update(ctx context.Context, comment *domains.TaskComment) error {
	//TODO implement me
	panic("implement me")
}

func (r *taskCommentRepo) Delete(ctx context.Context, id uint) error {
	//TODO implement me
	panic("implement me")
}

func (r *taskCommentRepo) GetListByTaskID(ctx context.Context, taskID uint, limit uint, offset uint) ([]domains.TaskComment, uint, error) {
	//TODO implement me
	panic("implement me")
}
