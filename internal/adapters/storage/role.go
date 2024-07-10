package storage

import (
	"context"
	"errors"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/mappers"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) ports.RoleRepository {
	return &roleRepository{
		db: db,
	}
}

var (
	ErrRoleAlreadyExists = errors.New("board already exists")
	ErrRoleNotFound      = errors.New("role not found")
)

func (r *roleRepository) Create(ctx context.Context, role *domains.Role) error {
	// Check if the role already exists
	var existingRole entities.Role
	err := r.db.WithContext(ctx).Model(&entities.Role{}).Where("name = ?", role.Name).First(&existingRole).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingRole.ID != 0 {
		return ErrRoleAlreadyExists // Custom error indicating role already exists
	}

	// Use transaction for creating the role
	return r.db.Transaction(func(tx *gorm.DB) error {
		entity := mappers.DomainToRoleEntity(role)

		if err := tx.WithContext(ctx).Create(entity).Error; err != nil {
			return err
		}
		role.ID = entity.ID // set the ID to the domain object after creation
		return nil
	})
}

func (r *roleRepository) GetByID(ctx context.Context, id uint) (*domains.Role, error) {
	var entity entities.Role
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {

		return nil, err
	}
	return mappers.RoleEntityToDomain(&entity), nil
}

func (r *roleRepository) GetAll(ctx context.Context) ([]domains.Role, error) {
	var entities []entities.Role
	if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return mappers.RoleEntitiesToDomain(entities), nil
}

func (r *roleRepository) Update(ctx context.Context, role *domains.Role) error {
	// Check if the role exists
	var existingRole entities.Role
	if err := r.db.WithContext(ctx).Model(&entities.Role{}).Where("id = ?", role.ID).First(&existingRole).Error; err != nil {
		return err
	}

	// Update fields
	existingRole.Name = role.Name
	existingRole.Description = role.Description

	// Save updated role
	return r.db.WithContext(ctx).Save(&existingRole).Error
}

func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Role{}, id).Error
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*domains.Role, error) {
	var entity entities.Role
	if err := r.db.WithContext(ctx).First(&entity, name).Error; err != nil {

		return nil, err
	}

	if entity.ID == 0 {
		return nil, ErrRoleNotFound
	}

	return mappers.RoleEntityToDomain(&entity), nil
}
