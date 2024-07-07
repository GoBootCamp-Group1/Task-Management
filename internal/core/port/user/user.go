package user

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
)

type Repo interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
