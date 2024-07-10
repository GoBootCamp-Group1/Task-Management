package ports

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
)

type RoleRepository interface {
	Create(ctx context.Context, role *domains.Role) error
	GetByID(ctx context.Context, id uint) (*domains.Role, error)
	GetAll(ctx context.Context) ([]domains.Role, error)
	Update(ctx context.Context, role *domains.Role) error
	Delete(ctx context.Context, id uint) error
}
