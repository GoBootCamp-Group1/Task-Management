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

type columnRepo struct {
	db *gorm.DB
}

var (
	ErrColumnAlreadyExists = errors.New("column already exists")
)

func NewColumnRepo(db *gorm.DB) port.ColumnRepo {
	return &columnRepo{
		db: db,
	}
}

func (r *columnRepo) Create(ctx context.Context, boardId uint, column *domain.Column) error {
	var existingColumn entities.Column
	var lastColumn entities.Column
	var lastPosition uint = 1

	err := r.db.WithContext(ctx).Model(&entities.Column{}).Where(&entities.Column{BoardId: boardId, Name: column.Name}).First(&existingColumn).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existingColumn.ID != 0 {
		return ErrColumnAlreadyExists
	}

	err = r.db.WithContext(ctx).Model(&entities.Column{}).Order("order_position DESC").First(&lastColumn).Error
	if err != nil {
		return err
	}
	if lastColumn.ID != 0 {
		lastPosition = lastColumn.OrderPosition + 1
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		entity := mappers.DomainToColumnEntity(column)
		entity.OrderPosition = lastPosition

		if err := tx.WithContext(ctx).Create(&entity).Error; err != nil {
			return err
		}

		column.ID = entity.ID
		return nil
	})

}

func (r *columnRepo) GetByID(ctx context.Context, id uint) (*domain.Column, error) {
	var column entities.Column
	err := r.db.WithContext(ctx).Model(&entities.Column{}).Where(&entities.Column{Model: gorm.Model{ID: id}}).First(&column).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mappers.ColumnEntityToDomain(&column), nil
}

func (r *columnRepo) GetAll(ctx context.Context, boardId uint) ([]*domain.Column, error) {
	var entitieColumns []entities.Column

	err := r.db.WithContext(ctx).Model(&entities.Column{}).Where(&entities.Column{BoardId: boardId}).Find(entitieColumns).Error
	if err != nil {
		return nil, nil
	}

	columns := make([]*domain.Column, len(entitieColumns))
	for i, c := range entitieColumns {
		columns[i] = mappers.ColumnEntityToDomain(&c)
	}

	return columns, nil
}

func (r *columnRepo) Update(ctx context.Context, column *domain.Column) error {
	var foundColumn *entities.Column
	err := r.db.WithContext(ctx).Model(&entities.Column{}).Where(&entities.Column{Model: gorm.Model{ID: column.ID}}).First(&foundColumn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

}
