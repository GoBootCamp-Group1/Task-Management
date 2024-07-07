package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

type Repo interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

var (
	ErrUserNotFound    = errors.New("User not found")
	ErrInvalidPassword = errors.New("Invalid user password")
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

func HashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
