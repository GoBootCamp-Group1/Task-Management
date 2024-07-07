package mappers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
	"gorm.io/gorm"
)

func UserEntityToDomain(entity *entities.User) *domain.User {
	return &domain.User{
		ID:       entity.ID,
		Name:     entity.Name,
		Email:    entity.Email,
		Password: entity.Password,
		Role:     domain.UserRole(entity.Role),
	}
}

func DomainToUserEntity(model *domain.User) *entities.User {
	return &entities.User{
		Model:    gorm.Model{ID: model.ID},
		Name:     model.Name,
		Email:    model.Email,
		Password: domain.HashPassword(model.Password),
	}
}
