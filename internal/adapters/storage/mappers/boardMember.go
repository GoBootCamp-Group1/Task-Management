package mappers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/fp"
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

func BoardMemberEntitiesToDomain(boardMemberEntities []entities.BoardMember) []domain.BoardMember {
	return fp.Map(boardMemberEntities, func(entity entities.BoardMember) domain.BoardMember {
		return *BoardMemberEntityToDomain(&entity)
	})
}

func BoardMemberDomainsToEntity(boardMemberDomains []domain.BoardMember) []entities.BoardMember {
	return fp.Map(boardMemberDomains, func(member domain.BoardMember) entities.BoardMember {
		return *DomainToBoardMemberEntity(&member)
	})
}
