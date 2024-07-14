package storage

import (
	"context"
	"errors"

	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/mappers"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) ports.UserRepo {
	return &userRepo{
		db: db,
	}
}

var (
	ErrUserAlreadyExists = "user already exists"
	ErrUserNotFound      = "user not found"
)

func (r *userRepo) Create(ctx context.Context, user *domains.User) error {
	var existingUser *entities.User
	err := r.db.WithContext(ctx).Model(&entities.User{}).Where("email = ?", user.Email).First(&existingUser).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if existingUser.ID != 0 {
		return fiber.NewError(fiber.StatusBadRequest, ErrUserAlreadyExists)
	}

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		entity := mappers.DomainToUserEntity(user)
		err := tx.WithContext(ctx).Create(&entity).Error
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id uint) (*domains.User, error) {
	var u entities.User

	err := r.db.WithContext(ctx).Model(&entities.User{}).Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if u.ID == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, ErrUserNotFound)
	}

	return mappers.UserEntityToDomain(&u), nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domains.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Model(&entities.User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, ErrUserNotFound)
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return mappers.UserEntityToDomain(&user), nil
}
