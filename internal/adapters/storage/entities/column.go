package entities

import "gorm.io/gorm"

type Column struct {
	gorm.Model
	BoardId       uint
	Name          string
	OrderPosition uint
	IsFinal       bool
	CreatedBy     uint
}
