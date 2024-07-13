package ports

import (
	"context"

	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
)

type ColumnRepo interface {
	Create(ctx context.Context, column *domains.Column) error
	GetByID(ctx context.Context, id uint) (*domains.Column, error)
	GetAll(ctx context.Context, boardId uint, limit int, offset int) ([]*domains.Column, error)
	Update(ctx context.Context, updateColumn *domains.ColumnUpdate) error
	Move(ctx context.Context, moveColumn *domains.ColumnMove) error
	Final(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
}
