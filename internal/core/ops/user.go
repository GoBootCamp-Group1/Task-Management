package ops

import (
	"context"
	user_model "github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/port"
)

type Ops struct {
	repo port.Repo
}

func NewOps(repo port.Repo) *Ops {
	return &Ops{
		repo: repo,
	}
}

func (o *Ops) Create(ctx context.Context, user *user_model.User) error {
	// validation
	return o.repo.Create(ctx, user)
}

func (o *Ops) GetUserByID(ctx context.Context, id uint) (*user_model.User, error) {
	return o.repo.GetByID(ctx, id)
}

func (o *Ops) GetUserByEmailAndPassword(ctx context.Context, email, password string) (*user_model.User, error) {
	user, err := o.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, user_model.ErrUserNotFound
	}

	if !user.PasswordIsValid(password) {
		return nil, user_model.ErrInvalidPassword
	}

	return user, nil
}
