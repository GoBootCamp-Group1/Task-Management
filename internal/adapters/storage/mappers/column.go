package mappers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
	"gorm.io/gorm"
)

func DomainToColumnEntity(column *domain.Column) *entities.Column {
	return &entities.Column{
		Model:         gorm.Model{ID: column.ID},
		CreatedBy:     column.CreatedBy,
		Name:          column.Name,
		IsFinal:       column.IsFinal,
		OrderPosition: column.OrderPosition,
		BoardID:       column.BoardID,
	}
}

func ColumnEntityToDomain(entity *entities.Column) *domain.Column {
	return &domain.Column{
		ID:            entity.ID,
		CreatedBy:     entity.CreatedBy,
		Name:          entity.Name,
		IsFinal:       entity.IsFinal,
		OrderPosition: entity.OrderPosition,
		BoardID:       entity.BoardID,
	}
}
