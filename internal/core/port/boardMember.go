package port

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
)

type BoardMemberRepo interface {
	Create(ctx context.Context, member *domain.BoardMember) error
	GetByID(ctx context.Context, id uint) (*domain.BoardMember, error)
	Update(ctx context.Context, member *domain.BoardMember) error
	Delete(ctx context.Context, id uint) error
}
