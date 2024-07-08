package storage

import (
	"context"
	"errors"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/mappers"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/port"
	"gorm.io/gorm"
)

type boardRepo struct {
	db *gorm.DB
}

func NewBoardRepo(db *gorm.DB) port.BoardRepo {
	return &boardRepo{
		db: db,
	}
}

func (r *boardRepo) Create(ctx context.Context, board *domain.Board) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		entity := mappers.DomainToBoardEntity(board)
		if err := tx.WithContext(ctx).Create(&entity).Error; err != nil {
			return err
		}
		board.ID = entity.ID // set the ID to the domain object after creation
		return nil
	})
}

func (r *boardRepo) GetByID(ctx context.Context, id uint) (*domain.Board, error) {
	var b entities.Board
	err := r.db.WithContext(ctx).Model(&entities.Board{}).Where("id = ?", id).First(&b).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mappers.BoardEntityToDomain(&b), nil
}

func (r *boardRepo) Update(ctx context.Context, board *domain.Board) error {
	entity := mappers.DomainToBoardEntity(board)
	return r.db.WithContext(ctx).Save(&entity).Error
}

func (r *boardRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Board{}, id).Error
}
