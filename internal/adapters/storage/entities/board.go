package entities

import (
	"gorm.io/gorm"
)

type Board struct {
	gorm.Model
	CreatedBy uint
	Name      string
	IsPrivate bool
}
