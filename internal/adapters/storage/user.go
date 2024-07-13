package storage

import (
	"context"
	"errors"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/mappers"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
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
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

func (r *userRepo) Create(ctx context.Context, user *domains.User) error {
	var existingUser *entities.User
	err := r.db.WithContext(ctx).Model(&entities.User{}).Where("email = ?", user.Email).First(&existingUser).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser.ID != 0 {
		return ErrUserAlreadyExists
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		entity := mappers.DomainToUserEntity(user)
		err := tx.WithContext(ctx).Create(&entity).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *userRepo) GetByID(ctx context.Context, id uint) (*domains.User, error) {
	var u entities.User

	err := r.db.WithContext(ctx).Model(&entities.User{}).Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}

	if u.ID == 0 {
		return nil, ErrUserNotFound
	}

	return mappers.UserEntityToDomain(&u), nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domains.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Model(&entities.User{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mappers.UserEntityToDomain(&user), nil
}
