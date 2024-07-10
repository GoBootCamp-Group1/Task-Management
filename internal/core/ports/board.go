package ports

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
)

type BoardRepo interface {
	Create(ctx context.Context, board *domains.Board) error
	GetByID(ctx context.Context, id uint) (*domains.Board, error)
	Update(ctx context.Context, board *domains.Board) error
	Delete(ctx context.Context, id uint) error
	GetAll(ctx context.Context) ([]domains.Board, error)
}
