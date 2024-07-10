package mappers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/fp"
	"gorm.io/gorm"
)

func DomainToRoleEntity(role *domain.Role) *entities.Role {
	return &entities.Role{
		Model: gorm.Model{
			ID: role.ID,
		},
		Name:        role.Name,
		Description: role.Description,
	}
}

func RoleEntityToDomain(entity *entities.Role) *domain.Role {
	return &domain.Role{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
	}
}

func RoleEntitiesToDomain(roleEntities []entities.Role) []domain.Role {
	return fp.Map(roleEntities, func(entity entities.Role) domain.Role {
		return *RoleEntityToDomain(&entity)
	})
}

func RoleDomainsToEntity(roleDomains []domain.Role) []entities.Role {
	return fp.Map(roleDomains, func(role domain.Role) entities.Role {
		return *DomainToRoleEntity(&role)
	})
}
