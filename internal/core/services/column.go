package services

import (
	"context"

	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/response"
)

type ColumnService struct {
	repo ports.ColumnRepo
}

func NewColumnService(repo ports.ColumnRepo) *ColumnService {
	return &ColumnService{repo: repo}
}

func (s *ColumnService) CreateColumn(ctx context.Context, column *domains.Column) error {
	return s.repo.Create(ctx, column)
}

func (s *ColumnService) GetColumnById(ctx context.Context, id uint) (*domains.Column, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ColumnService) GetAllColumns(ctx context.Context, boardId uint, limit int, offset int) (response.PaginateResponseFromService[[]*domains.Column], error) {
	return s.repo.GetAll(ctx, boardId, limit, offset)
}

func (s *ColumnService) Update(ctx context.Context, updateColumn *domains.ColumnUpdate) error {
	return s.repo.Update(ctx, updateColumn)
}

func (s *ColumnService) Move(ctx context.Context, moveColumn *domains.ColumnMove) error {
	return s.repo.Move(ctx, moveColumn)
}

func (s *ColumnService) Final(ctx context.Context, id uint) error {
	return s.repo.Final(ctx, id)
}

func (s *ColumnService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
