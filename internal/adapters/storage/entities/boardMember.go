package entities

import (
	"gorm.io/gorm"
)

type BoardMember struct {
	gorm.Model
	BoardID uint `gorm:"index"`
	UserID  uint `gorm:"index"`
	RoleID  uint `gorm:"index"`
}
