package domain

import (
	"gorm.io/gorm"
)

type BoardMember struct {
	gorm.Model
	BoardID uint
	UserID  uint
	Role    string
}
