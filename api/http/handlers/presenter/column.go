package presenter

import "github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"

type ColumnOutBoundPresenter struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	OrderPosition int    `json:"order_position"`
	IsFinal       bool   `json:"is_final"`
}

func NewColumnOutBoundPresenter(column *domains.Column) *ColumnOutBoundPresenter {
	return &ColumnOutBoundPresenter{
		ID:            column.ID,
		Name:          column.Name,
		OrderPosition: column.OrderPosition,
		IsFinal:       column.IsFinal,
	}
}
