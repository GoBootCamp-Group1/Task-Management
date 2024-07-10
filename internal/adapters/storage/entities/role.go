package entities

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255);not null;unique"`
	Description string `gorm:"type:text"`
}
