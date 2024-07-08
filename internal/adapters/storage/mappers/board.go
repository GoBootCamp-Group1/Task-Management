package mappers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
)

func DomainToBoardEntity(board *domain.Board) *entities.Board {
	return &entities.Board{
		ID:        board.ID,
		CreatedAt: board.CreatedAt,
		UpdatedAt: board.UpdatedAt,
		DeletedAt: board.DeletedAt,
		CreatedBy: board.CreatedBy,
		Name:      board.Name,
		IsPrivate: board.IsPrivate,
	}
}

func BoardEntityToDomain(entity *entities.Board) *domain.Board {
	return &domain.Board{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
		CreatedBy: entity.CreatedBy,
		Name:      entity.Name,
		IsPrivate: entity.IsPrivate,
	}
}
