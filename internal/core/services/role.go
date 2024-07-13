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

func (s *RoleService) InitRolesInDb(ctx context.Context) error {
	// check if roles exists in db already

	maintainerRoleEnum, _ := domains.ParseRole("Maintainer")
	editorRoleEnum, _ := domains.ParseRole("Editor")
	ownerRoleEnum, _ := domains.ParseRole("Owner")
	viewerRoleEnum, _ := domains.ParseRole("Viewer")

	maintainerRole := domains.Role{
		ID:          0,
		Name:        maintainerRoleEnum.String(),
		Description: "its maintainer role",
		Weight:      int(maintainerRoleEnum),
	}
	editorRole := domains.Role{
		ID:          0,
		Name:        editorRoleEnum.String(),
		Description: "its maintainer role",
		Weight:      int(editorRoleEnum),
	}
	ownerRole := domains.Role{
		ID:          0,
		Name:        ownerRoleEnum.String(),
		Description: "its maintainer role",
		Weight:      int(ownerRoleEnum),
	}
	viewerRole := domains.Role{
		ID:          0,
		Name:        viewerRoleEnum.String(),
		Description: "its maintainer role",
		Weight:      int(viewerRoleEnum),
	}

	roles := []domains.Role{maintainerRole, editorRole, ownerRole, viewerRole}
	for _, role := range roles {
		if err := s.CreateRole(ctx, &role); err != nil {
			// check if row already exists, then continue the for,
			//because if it exists we do not want to throw error or re-create
			return err
		}
	}

	return nil
}
