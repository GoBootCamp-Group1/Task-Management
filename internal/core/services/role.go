package services

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
)

type RoleService struct {
	roleRepo ports.RoleRepository
}

func NewRoleService(roleRepo ports.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (s *RoleService) CreateRole(ctx context.Context, role *domains.Role) error {
	return s.roleRepo.Create(ctx, role)
}
func (s *RoleService) UpdateRole(ctx context.Context, role *domains.Role) error {
	return s.roleRepo.Update(ctx, role)
}
func (s *RoleService) DeleteRole(ctx context.Context, id uint) error {
	return s.roleRepo.Delete(ctx, id)
}
func (s *RoleService) GetAllRoles(ctx context.Context) ([]domains.Role, error) {
	return s.roleRepo.GetAll(ctx)
}
func (s *RoleService) GetRoleById(ctx context.Context, id uint) (*domains.Role, error) {
	return s.roleRepo.GetByID(ctx, id)
}
