package services

import (
	"context"

	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

type BoardService struct {
	boardRepo       ports.BoardRepo
	boardMemberRepo ports.BoardMemberRepo
	userRepo        ports.UserRepo
	roleRepo        ports.RoleRepository
}

var (
	ErrUserIsAlreadyBoardMember = fiber.NewError(fiber.StatusBadRequest, "user is already a board member")
)

func NewBoardService(boardRepo ports.BoardRepo, boardMemberRepo ports.BoardMemberRepo, userRepo ports.UserRepo, roleRepo ports.RoleRepository) *BoardService {
	return &BoardService{boardRepo: boardRepo,
		boardMemberRepo: boardMemberRepo,
		userRepo:        userRepo,
		roleRepo:        roleRepo}
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

func (s *BoardService) GetRoleByUserIDAndBoardId(ctx context.Context, userID, boardID uint) (*domains.Role, error) {
	// Fetch board member details

	boardMember, err := s.boardMemberRepo.GetBoardMember(ctx, boardID, userID)
	if err != nil {
		return nil, err
	}

	// Fetch role details using roleRepo
	role, err := s.roleRepo.GetByID(ctx, boardMember.RoleID)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *BoardService) HasRequiredBoardAccess(ctx context.Context, roleW domains.RoleW, userID, boardID uint) (bool, error) {
	r, errGet := s.GetRoleByUserIDAndBoardId(ctx, userID, boardID)
	if errGet != nil {
		return false, errGet
	}
	roleWeight, errParse := domains.ParseRole(r.Name)
	if errParse != nil {
		return false, &fiber.Error{Code: fiber.StatusInternalServerError, Message: errParse.Error()}
	}
	if roleWeight > roleW {
		return false, &fiber.Error{Code: fiber.StatusForbidden, Message: "Access denied!"}
	}
	return true, nil
}

func (s *BoardService) InviteUserToBoard(ctx context.Context, userId, boardId uint, roleName string) error {
	// check board existence and get
	_, err := s.boardRepo.GetByID(ctx, boardId)
	if err != nil {
		return err
	}
	// check user existence and get
	_, err = s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}
	// check if user is member of board before
	_, err = s.boardMemberRepo.GetBoardMember(ctx, boardId, userId)
	if err == nil { // if err is nil, it means that the user is already board member
		return ErrUserIsAlreadyBoardMember
	}
	// check role existence and get
	role, err := s.roleRepo.GetByName(ctx, roleName)
	if err != nil {
		return err
	}
	// add row to board access and board member
	boardMember := &domains.BoardMember{
		BoardID: boardId,
		UserID:  userId,
		RoleID:  role.ID,
	}
	err = s.boardMemberRepo.Create(ctx, boardMember)
	if err != nil {
		return err
	}

	return nil
}

func (s *BoardService) RemoveUserFromBoard(ctx context.Context, userId, boardId uint) error {
	// check board existence and get
	_, err := s.boardRepo.GetByID(ctx, boardId)
	if err != nil {
		return err
	}
	// check user existence and get
	_, err = s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}
	// check if user is member of board before
	boardMember, err := s.boardMemberRepo.GetBoardMember(ctx, boardId, userId)
	if err != nil {
		return err
	}

	// remove user's board membership
	if err = s.boardMemberRepo.Delete(ctx, boardMember.ID); err != nil {
		return err
	}

	return nil
}

func (s *BoardService) ChangeUserRole(ctx context.Context, userId, boardId uint, roleName string) error {
	// check board existence and get
	_, err := s.boardRepo.GetByID(ctx, boardId)
	if err != nil {
		return err
	}
	// check user existence and get
	_, err = s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}
	// check if user is member of board before
	boardMember, err := s.boardMemberRepo.GetBoardMember(ctx, boardId, userId)
	if err != nil {
		return err
	}
	// check role existence and get
	role, err := s.roleRepo.GetByName(ctx, roleName)
	if err != nil {
		return err
	}

	boardMember.RoleID = role.ID
	// change user role in board
	if err = s.boardMemberRepo.Update(ctx, boardMember); err != nil {
		return err
	}

	return nil
}
