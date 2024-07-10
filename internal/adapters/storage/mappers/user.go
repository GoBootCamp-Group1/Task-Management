package mappers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/utils"
	"gorm.io/gorm"
)

func UserEntityToDomain(entity *entities.User) *domains.User {
	return &domains.User{
		ID:       entity.ID,
		Name:     entity.Name,
		Email:    entity.Email,
		Password: entity.Password,
		Role:     domains.UserRole(entity.Role),
	}
}

func DomainToUserEntity(model *domains.User) *entities.User {
	return &entities.User{
		Model:    gorm.Model{ID: model.ID},
		Name:     model.Name,
		Email:    model.Email,
		Password: utils.HashPassword(model.Password),
	}
}
