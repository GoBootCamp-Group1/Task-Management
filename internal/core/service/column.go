package service

import (
	"context"

	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/port"
)

type ColumnService struct {
	repo port.ColumnRepo
}

func NewColumnService(repo port.ColumnRepo) *ColumnService {
	return &ColumnService{repo: repo}
}

func (s *ColumnService) CreateColumn(ctx context.Context, column *domain.Column) error {
	return s.repo.Create(ctx, column)
}

func (s *ColumnService) GetColumnById(ctx context.Context, id uint) (*domain.Column, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ColumnService) GetAllColumns(ctx context.Context, boardId uint) ([]*domain.Column, error) {
	return s.repo.GetAll(ctx, boardId)
}

func (s *ColumnService) Update(ctx context.Context, updateColumn *domain.ColumnUpdate) error {
	return s.repo.Update(ctx, updateColumn)
}

func (s *ColumnService) Move(ctx context.Context, moveColumn *domain.ColumnMove) error {
	return s.repo.Move(ctx, moveColumn)
}

func (s *ColumnService) Final(ctx context.Context, id uint) error {
	return s.repo.Final(ctx, id)
}

func (s *ColumnService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
