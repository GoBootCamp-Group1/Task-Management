package mappers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain/user"
	"gorm.io/gorm"
)

func UserEntityToDomain(entity *entities.User) *user.User {
	return &user.User{
		ID:       entity.ID,
		Name:     entity.Name,
		Email:    entity.Email,
		Password: entity.Password,
		Role:     user.UserRole(entity.Role),
	}
}

func DomainToUserEntity(model *user.User) *entities.User {
	return &entities.User{
		Model:    gorm.Model{ID: model.ID},
		Name:     model.Name,
		Email:    model.Email,
		Password: user.HashPassword(model.Password),
	}
}
