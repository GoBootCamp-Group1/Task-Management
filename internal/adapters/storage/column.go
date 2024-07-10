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

type columnRepo struct {
	db *gorm.DB
}

var (
	ErrColumnAlreadyExists     = errors.New("column already exists")
	ErrOrderPositionOutOfRange = errors.New("order position is out of range")
)

func NewColumnRepo(db *gorm.DB) ports.ColumnRepo {
	return &columnRepo{
		db: db,
	}
}

func (r *columnRepo) Create(ctx context.Context, column *domains.Column) error {
	var existingColumn entities.Column
	var lastColumn entities.Column
	var lastPosition int = 1

	err := r.db.WithContext(ctx).Model(&entities.Column{}).Where(&entities.Column{BoardID: column.BoardID, Name: column.Name}).First(&existingColumn).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existingColumn.ID != 0 {
		return ErrColumnAlreadyExists
	}

	err = r.db.WithContext(ctx).Model(&entities.Column{}).Order("order_position DESC").First(&lastColumn).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

func (r *columnRepo) GetByID(ctx context.Context, id uint) (*domains.Column, error) {
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

func (r *columnRepo) GetAll(ctx context.Context, boardId uint) ([]*domains.Column, error) {
	var entitieColumns []entities.Column

	err := r.db.WithContext(ctx).Model(&entities.Column{}).Where(&entities.Column{BoardID: boardId}).Order("order_position ASC").Find(&entitieColumns).Error
	if err != nil {
		return nil, err
	}

	columns := make([]*domains.Column, len(entitieColumns))
	for i, c := range entitieColumns {
		columns[i] = mappers.ColumnEntityToDomain(&c)
	}

	return columns, nil
}

// only name can update
func (r *columnRepo) Update(ctx context.Context, updateColumn *domains.ColumnUpdate) error {
	var foundColumn *entities.Column
	err := r.db.WithContext(ctx).Model(&entities.Column{}).Where(&entities.Column{Model: gorm.Model{ID: updateColumn.ID}}).First(&foundColumn).Error
	if err != nil {
		return err
	}

	foundColumn.Name = updateColumn.Name

	return r.db.WithContext(ctx).Save(&foundColumn).Error
}

func (r *columnRepo) Move(ctx context.Context, moveColumn *domains.ColumnMove) error {
	var foundColumn *entities.Column
	var lastColumn entities.Column
	err := r.db.WithContext(ctx).Model(&entities.Column{}).Where(&entities.Column{Model: gorm.Model{ID: moveColumn.ID}}).First(&foundColumn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	if foundColumn.OrderPosition == moveColumn.OrderPosition {
		return nil
	}

	err = r.db.WithContext(ctx).Model(&entities.Column{}).Order("order_position DESC").First(&lastColumn).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if moveColumn.OrderPosition < 1 || moveColumn.OrderPosition > lastColumn.OrderPosition {
		return ErrOrderPositionOutOfRange
	}

	var allColumns []*entities.Column
	var condition string
	var unit int

	if foundColumn.OrderPosition > moveColumn.OrderPosition {
		condition = "board_id = ? AND order_position >= ? AND order_position < ?"
		unit = 1
	} else {
		condition = "board_id = ? AND order_position <= ? AND order_position > ?"
		unit = -1
	}

	err = r.db.WithContext(ctx).Model(&entities.Column{}).Where(condition, foundColumn.BoardID, moveColumn.OrderPosition, foundColumn.OrderPosition).Order("order_position ASC").Find(&allColumns).Error
	if err != nil {
		return err
	}
	for _, col := range allColumns {
		col.OrderPosition += unit
		err = r.db.WithContext(ctx).Save(&col).Error
		if err != nil {
			return err
		}
	}

	foundColumn.OrderPosition = moveColumn.OrderPosition
	return r.db.WithContext(ctx).Save(&foundColumn).Error
}

func (r *columnRepo) Final(ctx context.Context, id uint) error {
	var foundColumn *entities.Column
	err := r.db.WithContext(ctx).Model(&entities.Column{}).Where(&entities.Column{Model: gorm.Model{ID: id}}).First(&foundColumn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	if foundColumn.IsFinal {
		return nil
	}

	err = r.db.Model(&entities.Column{}).Where(&entities.Column{BoardID: foundColumn.BoardID}).Update("is_final", false).Error
	if err != nil {
		return nil
	}

	foundColumn.IsFinal = true

	return r.db.WithContext(ctx).Save(&foundColumn).Error
}

func (r *columnRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Column{}, id).Error
}
