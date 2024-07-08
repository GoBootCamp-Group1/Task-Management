package mappers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
)

func DomainToBoardMemberEntity(member *domain.BoardMember) *entities.BoardMember {
	return &entities.BoardMember{
		Model:   member.Model,
		BoardID: member.BoardID,
		UserID:  member.UserID,
		Role:    member.Role,
	}
}

func BoardMemberEntityToDomain(entity *entities.BoardMember) *domain.BoardMember {
	return &domain.BoardMember{
		Model:   entity.Model,
		BoardID: entity.BoardID,
		UserID:  entity.UserID,
		Role:    entity.Role,
	}
}
