package service

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/port"
)

type BoardService struct {
	boardRepo       port.BoardRepo
	boardMemberRepo port.BoardMemberRepo
	userRepo        port.UserRepo
}

func NewBoardService(repo port.BoardRepo) *BoardService {
	return &BoardService{boardRepo: repo}
}

func (s *BoardService) CreateBoard(ctx context.Context, board *domain.Board) error {
	return s.boardRepo.Create(ctx, board)
}

func (s *BoardService) GetBoardByID(ctx context.Context, id uint) (*domain.Board, error) {
	return s.boardRepo.GetByID(ctx, id)
}

func (s *BoardService) UpdateBoard(ctx context.Context, board *domain.Board) error {
	return s.boardRepo.Update(ctx, board)
}

func (s *BoardService) DeleteBoard(ctx context.Context, id uint) error {
	return s.boardRepo.Delete(ctx, id)
}

func (s *BoardService) GetAllBoards(ctx context.Context) ([]*domain.Board, error) {
	return s.boardRepo.GetAll(ctx)
}

func (s *BoardService) CreateBoardMember(ctx context.Context, boardMember *domain.BoardMember) error {
	return s.boardMemberRepo.Create(ctx, boardMember)
}
func (s *BoardService) GetBoardMembersByBoardId(ctx context.Context, id uint) ([]*domain.User, error) {
	var users []*domain.User

	boardMember, err := s.boardMemberRepo.GetBoardMembers(ctx, id)
	if err != nil {
		return users, err
	}
	for _, member := range boardMember {
		user, err := s.userRepo.GetByID(ctx, member.UserID)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}
