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

type boardRepo struct {
	db *gorm.DB
}

func NewBoardRepo(db *gorm.DB) ports.BoardRepo {
	return &boardRepo{
		db: db,
	}
}

var (
	ErrBoardAlreadyExists = "Board already exists"
)

func (r *boardRepo) Create(ctx context.Context, board *domains.Board) error {
	var existingBoard entities.Board
	err := r.db.WithContext(ctx).Model(&entities.Board{}).Where("name = ? AND created_by = ?", board.Name, board.CreatedBy).First(&existingBoard).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if existingBoard.ID != 0 {
		return fiber.NewError(fiber.StatusBadRequest, ErrBoardAlreadyExists)
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		entity := mappers.DomainToBoardEntity(board)

		if err := tx.WithContext(ctx).Create(&entity).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		board.ID = entity.ID // set the ID to the domain object after creation
		return nil
	})
}

func (r *boardRepo) GetByID(ctx context.Context, id uint) (*domains.Board, error) {
	var b entities.Board
	err := r.db.WithContext(ctx).Model(&entities.Board{}).Where("id = ?", id).First(&b).Error
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return mappers.BoardEntityToDomain(&b), nil
}

func (r *boardRepo) Update(ctx context.Context, board *domains.Board) error {
	var existingBoard *entities.Board
	if err := r.db.WithContext(ctx).Model(&entities.Board{}).Where("id = ?", board.ID).First(&existingBoard).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	existingBoard.Name = board.Name
	existingBoard.IsPrivate = board.IsPrivate

	if err := r.db.WithContext(ctx).Save(&existingBoard).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (r *boardRepo) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entities.Board{}, id).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *boardRepo) GetAll(ctx context.Context) ([]domains.Board, error) {
	var boards []entities.Board
	err := r.db.WithContext(ctx).Find(&boards).Error
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return mappers.BoardEntitiesToDomain(boards), nil
}
