package port

import (
	"context"

	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
)

type ColumnRepo interface {
	Create(ctx context.Context, boardId uint, column *domain.Column) error
	GetByID(ctx context.Context, id uint) (*domain.Column, error)
	GetAll(ctx context.Context, boardId uint) ([]*domain.Column, error)
	Update(ctx context.Context, column *domain.Column) error
	Delete(ctx context.Context, id uint) error
	Move(ctx context.Context, column []*domain.Column) error
}
