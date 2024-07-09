package ports

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
)

type UserRepo interface {
	Create(ctx context.Context, user *domains.User) error
	GetByID(ctx context.Context, id uint) (*domains.User, error)
	GetByEmail(ctx context.Context, email string) (*domains.User, error)
}
