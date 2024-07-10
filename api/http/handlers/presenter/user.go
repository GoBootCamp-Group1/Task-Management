package presenter

import "github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"

type UserPresenter struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUserPresenter(user *domains.User) *UserPresenter {
	return &UserPresenter{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
