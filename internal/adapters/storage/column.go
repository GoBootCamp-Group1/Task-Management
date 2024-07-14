package storage

import (
	"context"
	"errors"

	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/mappers"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/response"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type columnRepo struct {
	db *gorm.DB
}

var (
	ErrColumnAlreadyExists     = "column already exists"
	ErrOrderPositionOutOfRange = "order position is out of range"
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
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if existingColumn.ID != 0 {
		return fiber.NewError(fiber.StatusBadRequest, ErrColumnAlreadyExists)
	}

	err = r.db.WithContext(ctx).Model(&entities.Column{}).Order("order_position DESC").First(&lastColumn).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if lastColumn.ID != 0 {
		lastPosition = lastColumn.OrderPosition + 1
	}

	if column.IsFinal {
		err = r.db.Model(&entities.Column{}).Where(&entities.Column{BoardID: column.BoardID}).Update("is_final", false).Error
		if err != nil {
			return nil
		}
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		entity := mappers.DomainToColumnEntity(column)
		entity.OrderPosition = lastPosition

		if err := tx.WithContext(ctx).Create(&entity).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
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
			return nil, fiber.NewError(fiber.StatusNotFound, "Column not found")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return mappers.ColumnEntityToDomain(&column), nil
}

func (r *columnRepo) GetAll(ctx context.Context, boardId uint, page int, pageSize int) (response.PaginateResponseFromService[[]*domains.Column], error) {
	var columnEntities []entities.Column

	query := r.db.WithContext(ctx).Model(&entities.Column{}).Where(&entities.Column{BoardID: boardId}).Order("order_position ASC")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return response.PaginateResponseFromService[[]*domains.Column]{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if pageSize > 0 {
		query = query.Offset((page * pageSize) - pageSize)
	}

	if page > 0 {
		query = query.Limit(pageSize)
	}

	err := query.Debug().Find(&columnEntities).Error

	if err != nil {
		return response.PaginateResponseFromService[[]*domains.Column]{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	columns := make([]*domains.Column, len(columnEntities))
	for i, c := range columnEntities {
		columns[i] = mappers.ColumnEntityToDomain(&c)
	}

	return response.PaginateResponseFromService[[]*domains.Column]{
		Data:     columns,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

// only name can update
func (r *columnRepo) Update(ctx context.Context, updateColumn *domains.ColumnUpdate) error {
	var foundColumn *entities.Column
	err := r.db.WithContext(ctx).Model(&entities.Column{}).Where(&entities.Column{Model: gorm.Model{ID: updateColumn.ID}}).First(&foundColumn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Column not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	foundColumn.Name = updateColumn.Name

	return r.db.WithContext(ctx).Save(&foundColumn).Error
}

func (r *columnRepo) Move(ctx context.Context, moveColumn *domains.ColumnMove) error {
	var foundColumn *entities.Column
	var lastColumn entities.Column
	err := r.db.WithContext(ctx).Model(&entities.Column{}).Where("id = ?", moveColumn.ID).First(&foundColumn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Column not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if foundColumn.OrderPosition == moveColumn.OrderPosition {
		return nil
	}

	err = r.db.WithContext(ctx).Model(&entities.Column{}).Order("order_position DESC").First(&lastColumn).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if moveColumn.OrderPosition < 1 || moveColumn.OrderPosition > lastColumn.OrderPosition {
		return fiber.NewError(fiber.StatusBadRequest, ErrOrderPositionOutOfRange)
	}

	var condition string
	var unit int

	if foundColumn.OrderPosition > moveColumn.OrderPosition {
		condition = "board_id = ? AND order_position >= ? AND order_position < ?"
		unit = 1
	} else {
		condition = "board_id = ? AND order_position <= ? AND order_position > ?"
		unit = -1
	}

	err = r.db.WithContext(ctx).Model(&entities.Column{}).Where(condition, foundColumn.BoardID, moveColumn.OrderPosition, foundColumn.OrderPosition).Updates(map[string]interface{}{"order_position": gorm.Expr("order_position + ?", unit)}).Error
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	foundColumn.OrderPosition = moveColumn.OrderPosition
	err = r.db.WithContext(ctx).Save(&foundColumn).Error
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *columnRepo) Final(ctx context.Context, id uint) error {
	var foundColumn *entities.Column
	err := r.db.WithContext(ctx).Model(&entities.Column{}).Where(&entities.Column{Model: gorm.Model{ID: id}}).First(&foundColumn).Error
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if foundColumn.IsFinal {
		return nil
	}

	err = r.db.Model(&entities.Column{}).Where(&entities.Column{BoardID: foundColumn.BoardID}).Update("is_final", false).Error
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	foundColumn.IsFinal = true

	err = r.db.WithContext(ctx).Save(&foundColumn).Error
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *columnRepo) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entities.Column{}, id).Error
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}
