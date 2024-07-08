package entities

import "time"

type Board struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	CreatedBy uint
	Name      string
	IsPrivate bool
}
