package service

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/port"
)

type BoardService struct {
	repo port.BoardRepo
}

func NewBoardService(repo port.BoardRepo) *BoardService {
	return &BoardService{repo: repo}
}

func (s *BoardService) CreateBoard(ctx context.Context, board *domain.Board) error {
	return s.repo.Create(ctx, board)
}

func (s *BoardService) GetBoardByID(ctx context.Context, id uint) (*domain.Board, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BoardService) UpdateBoard(ctx context.Context, board *domain.Board) error {
	return s.repo.Update(ctx, board)
}

func (s *BoardService) DeleteBoard(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
