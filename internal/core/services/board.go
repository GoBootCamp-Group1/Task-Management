package services

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
)

type BoardService struct {
	repo ports.BoardRepo
}

func NewBoardService(repo ports.BoardRepo) *BoardService {
	return &BoardService{repo: repo}
}

func (s *BoardService) CreateBoard(ctx context.Context, board *domains.Board) error {
	return s.repo.Create(ctx, board)
}

func (s *BoardService) GetBoardByID(ctx context.Context, id uint) (*domains.Board, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BoardService) UpdateBoard(ctx context.Context, board *domains.Board) error {
	return s.repo.Update(ctx, board)
}

func (s *BoardService) DeleteBoard(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
