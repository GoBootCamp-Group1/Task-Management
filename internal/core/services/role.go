package services

import (
	"context"
	"errors"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
	"log"
	"time"
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

	maintainerRoleEnum, _ := domains.ParseRole("Maintainer")
	editorRoleEnum, _ := domains.ParseRole("Editor")
	ownerRoleEnum, _ := domains.ParseRole("Owner")
	viewerRoleEnum, _ := domains.ParseRole("Viewer")

	roles := []domains.Role{
		{
			Name:        maintainerRoleEnum.String(),
			Description: "its maintainer role",
			Weight:      int(maintainerRoleEnum),
		},
		{
			Name:        editorRoleEnum.String(),
			Description: "its editor role",
			Weight:      int(editorRoleEnum),
		},
		{
			Name:        ownerRoleEnum.String(),
			Description: "its owner role",
			Weight:      int(ownerRoleEnum),
		},
		{
			Name:        viewerRoleEnum.String(),
			Description: "its viewer role",
			Weight:      int(viewerRoleEnum),
		},
	}
	var lastErr error
	for _, role := range roles {
		err := createRoleWithRetry(s, role, ctx)
		if err != nil {
			lastErr = err
			log.Printf("Failed to create role %s: %v", role.Name, err)
		}
	}
	return lastErr
}
func createRoleWithRetry(s *RoleService, role domains.Role, ctx context.Context) error {

	const maxRetries = 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		err := s.CreateRole(ctx, &role)
		if err != nil {
			if errors.Is(err, storage.ErrRoleAlreadyExists) {
				// Role already exists, no need to retry
				return nil
			}
			// Log the error and prepare for retry
			log.Printf("Attempt %d: Failed to create role %s: %v", i+1, role.Name, err)
			lastErr = err
			time.Sleep(time.Second * time.Duration(i+1)) // Exponential backoff
		} else {
			// Role created successfully
			return nil
		}
	}

	// If we exit the loop, it means all retries failed
	return lastErr
}
