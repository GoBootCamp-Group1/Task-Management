package mappers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
)

func DomainToBoardMemberEntity(member *domains.BoardMember) *entities.BoardMember {
	return &entities.BoardMember{
		Model:   member.Model,
		BoardID: member.BoardID,
		UserID:  member.UserID,
		Role:    member.Role,
	}
}

func BoardMemberEntityToDomain(entity *entities.BoardMember) *domains.BoardMember {
	return &domains.BoardMember{
		Model:   entity.Model,
		BoardID: entity.BoardID,
		UserID:  entity.UserID,
		Role:    entity.Role,
	}
}
