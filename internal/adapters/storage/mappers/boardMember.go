package mappers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
	"gorm.io/gorm"
)

func DomainToBoardMemberEntity(member *domain.BoardMember) *entities.BoardMember {
	return &entities.BoardMember{
		Model: gorm.Model{
			ID: member.ID, // Set the ID from the domain model
		},
		BoardID: member.BoardID,
		UserID:  member.UserID,
		RoleID:  member.RoleID,
	}
}

func BoardMemberEntityToDomain(entity *entities.BoardMember) *domain.BoardMember {
	return &domain.BoardMember{
		ID:      entity.ID,
		BoardID: entity.BoardID,
		UserID:  entity.UserID,
		RoleID:  entity.RoleID,
	}
}
