package mappers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/fp"
	"gorm.io/gorm"
)

func DomainToRoleEntity(role *domains.Role) *entities.Role {
	return &entities.Role{
		Model: gorm.Model{
			ID: role.ID,
		},
		Name:        role.Name,
		Description: role.Description,
		Weight:      role.Weight,
	}
}

func RoleEntityToDomain(entity *entities.Role) *domains.Role {
	return &domains.Role{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Weight:      entity.Weight,
	}
}

func RoleEntitiesToDomain(roleEntities []entities.Role) []domains.Role {
	return fp.Map(roleEntities, func(entity entities.Role) domains.Role {
		return *RoleEntityToDomain(&entity)
	})
}

func RoleDomainsToEntity(roleDomains []domains.Role) []entities.Role {
	return fp.Map(roleDomains, func(role domains.Role) entities.Role {
		return *DomainToRoleEntity(&role)
	})
}
