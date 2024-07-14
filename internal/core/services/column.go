package services

import (
	"context"

	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type ColumnService struct {
	repo         ports.ColumnRepo
	boardService *BoardService
}

func NewColumnService(repo ports.ColumnRepo, boardService *BoardService) *ColumnService {
	return &ColumnService{
		repo:         repo,
		boardService: boardService,
	}
}

func (s *ColumnService) CreateColumn(ctx context.Context, column *domains.Column) error {
	hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Maintainer, column.CreatedBy, column.BoardID)
	if !hasAccess {
		return &fiber.Error{Code: fiber.StatusForbidden, Message: "Access denied"}
	}
	return s.repo.Create(ctx, column)
}

func (s *ColumnService) GetColumnById(ctx context.Context, userID uint, id uint) (*domains.Column, error) {
	column, errColumn := s.repo.GetByID(ctx, id)

	if errColumn != nil {
		return nil, errColumn
	}

	if column.Board.IsPrivate {
		//check permissions -> only board members can see task
		hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Viewer, userID, column.BoardID)
		if !hasAccess {
			return nil, &fiber.Error{Code: fiber.StatusForbidden, Message: "Access denied"}
		}
	}

	return column, nil
}

func (s *ColumnService) GetAllColumns(ctx context.Context, userID uint, boardID uint, limit int, offset int) (response.PaginateResponseFromService[[]*domains.Column], error) {
	hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Viewer, userID, boardID)
	if !hasAccess {
		return response.PaginateResponseFromService[[]*domains.Column]{}, &fiber.Error{Code: fiber.StatusForbidden, Message: "Access denied"}
	}
	return s.repo.GetAll(ctx, boardID, limit, offset)
}

func (s *ColumnService) Update(ctx context.Context, boardID uint, userID uint, updateColumn *domains.ColumnUpdate) error {
	hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Editor, userID, boardID)
	if !hasAccess {
		return &fiber.Error{Code: fiber.StatusForbidden, Message: "Access denied"}
	}
	return s.repo.Update(ctx, updateColumn)
}

func (s *ColumnService) Move(ctx context.Context, boardID uint, userID uint, moveColumn *domains.ColumnMove) error {
	hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Editor, userID, boardID)
	if !hasAccess {
		return &fiber.Error{Code: fiber.StatusForbidden, Message: "Access denied"}
	}
	return s.repo.Move(ctx, moveColumn)
}

func (s *ColumnService) Final(ctx context.Context, boardID uint, userID uint, id uint) error {
	hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Editor, userID, boardID)
	if !hasAccess {
		return &fiber.Error{Code: fiber.StatusForbidden, Message: "Access denied"}
	}
	return s.repo.Final(ctx, id)
}

func (s *ColumnService) Delete(ctx context.Context, boardID uint, userID uint, id uint) error {
	hasAccess, _ := s.boardService.HasRequiredBoardAccess(ctx, domains.Editor, userID, boardID)
	if !hasAccess {
		return &fiber.Error{Code: fiber.StatusForbidden, Message: "Access denied"}
	}
	return s.repo.Delete(ctx, id)
}
