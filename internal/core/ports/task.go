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
}
