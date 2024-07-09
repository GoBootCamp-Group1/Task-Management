package ports

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
)

type BoardMemberRepo interface {
	Create(ctx context.Context, member *domains.BoardMember) error
	GetByID(ctx context.Context, id uint) (*domains.BoardMember, error)
	Update(ctx context.Context, member *domains.BoardMember) error
	Delete(ctx context.Context, id uint) error
}
