package mappers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/fp"
	"gorm.io/gorm"
)

func DomainToBoardEntity(board *domains.Board) *entities.Board {
	return &entities.Board{
		Model:     gorm.Model{ID: board.ID},
		CreatedBy: board.CreatedBy,
		Name:      board.Name,
		IsPrivate: board.IsPrivate,
	}
}

func BoardEntityToDomain(entity *entities.Board) *domains.Board {
	return &domains.Board{
		ID:        entity.ID,
		CreatedBy: entity.CreatedBy,
		Name:      entity.Name,
		IsPrivate: entity.IsPrivate,
	}
}

func BoardEntitiesToDomain(boardEntities []entities.Board) []domain.Board {
	return fp.Map(boardEntities, func(entity entities.Board) domain.Board {
		return *BoardEntityToDomain(&entity)
	})
}

func BoardDomainsToEntity(boardDomains []domains.Board) []entities.Board {
	return fp.Map(boardDomains, func(member domains.Board) entities.Board {
		return *DomainToBoardEntity(&member)
	})
}
