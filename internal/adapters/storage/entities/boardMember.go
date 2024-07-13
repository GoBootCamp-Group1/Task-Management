package entities

import (
	"gorm.io/gorm"
)

type BoardMember struct {
	gorm.Model
	BoardID uint
	UserID  uint
	RoleID  uint
}

func (BoardMember) TableName() string {
	return "board_users"
}
