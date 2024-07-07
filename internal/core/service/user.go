package service

import (
	"context"
	user2 "github.com/GoBootCamp-Group1/Task-Management/internal/core/domain/user"
)

type UserService struct {
	userOps *user2.Ops
}

func NewUserService(userOps *user2.Ops) *UserService {
	return &UserService{
		userOps: userOps,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *user2.User) error {
	return s.userOps.Create(ctx, user)
}
