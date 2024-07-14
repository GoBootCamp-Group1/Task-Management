package entities

import "gorm.io/gorm"

type Column struct {
	gorm.Model
	BoardID       uint
	Name          string
	OrderPosition int
	IsFinal       bool
	CreatedBy     uint

	Board *Board
}
