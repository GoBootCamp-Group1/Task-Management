package domain

import (
	"crypto/sha256"
	"errors"
	"fmt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid user password")
)

type UserRole uint8

func (ur UserRole) String() string {
	switch ur {
	case UserRoleUser:
		return "user"
	case UserRoleAdmin:
		return "admin"
	default:
		return "unknown"
	}
}

const (
	UserRoleUser UserRole = iota + 1
	UserRoleAdmin
)

type User struct {
	ID       uint
	Name     string
	Email    string
	Password string
	Role     UserRole
}

func (u *User) ValidatePassword() error {
	return nil
}

func (u *User) PasswordIsValid(pass string) bool {
	h := sha256.New()
	h.Write([]byte(pass))
	passSha256 := h.Sum(nil)
	return fmt.Sprintf("%x", passSha256) == u.Password
}
