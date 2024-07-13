package entities

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	CreatedBy uint

	BoardID       uint
	ParentID      *uint
	AssigneeID    *uint
	ColumnID      uint
	OrderPosition int
	Name          string
	Description   string
	StartDateTime *time.Time `gorm:"column:start_datetime"`
	EndDateTime   *time.Time `gorm:"column:end_datetime"`
	StoryPoint    int
	Additional    json.RawMessage

	Board    Board  `gorm:"foreignKey:BoardID"`
	Creator  User   `gorm:"foreignKey:CreatedBy"`
	Column   Column `gorm:"foreignKey:ColumnID"`
	Parent   *Task  `gorm:"foreignKey:ParentID"`
	Assignee *User  `gorm:"foreignKey:AssigneeID"`
}
