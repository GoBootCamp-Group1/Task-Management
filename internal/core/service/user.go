package service

import (
	"context"
	user2 "github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ops"
)

type UserService struct {
	userOps *ops.Ops
}

func NewUserService(userOps *ops.Ops) *UserService {
	return &UserService{
		userOps: userOps,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *user2.User) error {
	return s.userOps.Create(ctx, user)
}
