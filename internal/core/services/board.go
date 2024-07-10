package services

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
)

type BoardService struct {
	boardRepo       ports.BoardRepo
	boardMemberRepo ports.BoardMemberRepo
	userRepo        ports.UserRepo
}

func NewBoardService(boardRepo ports.BoardRepo, boardMemberRepo ports.BoardMemberRepo, userRepo ports.UserRepo) *BoardService {
	return &BoardService{boardRepo: boardRepo,
		boardMemberRepo: boardMemberRepo,
		userRepo:        userRepo}
}

func (s *BoardService) CreateBoard(ctx context.Context, board *domains.Board) error {
	return s.boardRepo.Create(ctx, board)
}

func (s *BoardService) GetBoardByID(ctx context.Context, id uint) (*domains.Board, error) {
	return s.boardRepo.GetByID(ctx, id)
}

func (s *BoardService) UpdateBoard(ctx context.Context, board *domains.Board) error {
	return s.boardRepo.Update(ctx, board)
}

func (s *BoardService) DeleteBoard(ctx context.Context, id uint) error {
	return s.boardRepo.Delete(ctx, id)
}

func (s *BoardService) GetAllBoards(ctx context.Context) ([]domains.Board, error) {
	return s.boardRepo.GetAll(ctx)
}

func (s *BoardService) CreateBoardMember(ctx context.Context, boardMember *domains.BoardMember) error {
	return s.boardMemberRepo.Create(ctx, boardMember)
}
func (s *BoardService) GetBoardMembersByBoardId(ctx context.Context, boardId uint) ([]*domains.User, error) {
	var users []*domains.User

	boardMember, err := s.boardMemberRepo.GetBoardMembers(ctx, boardId)
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