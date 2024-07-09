package port

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
)

type BoardRepo interface {
	Create(ctx context.Context, board *domain.Board) error
	GetByID(ctx context.Context, id uint) (*domain.Board, error)
	Update(ctx context.Context, board *domain.Board) error
	Delete(ctx context.Context, id uint) error
	GetAll(ctx context.Context) ([]*domain.Board, error)
}
