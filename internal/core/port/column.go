package port

import (
	"context"

	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
)

type ColumnRepo interface {
	Create(ctx context.Context, column *domain.Column) error
	GetByID(ctx context.Context, id uint) (*domain.Column, error)
	GetAll(ctx context.Context, boardId uint) ([]*domain.Column, error)
	Update(ctx context.Context, updateColumn *domain.ColumnUpdate) error
	Move(ctx context.Context, moveColumn *domain.ColumnMove) error
	Final(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
}
