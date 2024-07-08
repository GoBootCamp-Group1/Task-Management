package service

import (
	"context"
	user_model "github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/port"
)

type UserService struct {
	repo port.UserRepo
}

func NewUserService(repo port.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *user_model.User) error {
	return s.repo.Create(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*user_model.User, error) {
	return s.repo.GetByID(ctx, id)
}
