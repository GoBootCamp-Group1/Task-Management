package entities

import (
	"encoding/json"
	"github.com/google/uuid"
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

type TaskChild struct {
	gorm.Model
	ColumnID      uint
	OrderPosition int
	Name          string
	Description   string
	ColumnName    string `gorm:"column:column_name"`
	ColumnIsFinal bool   `gorm:"column:is_final"`
}

type TaskComment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    uint
	TaskID    uint
	Comment   string
	Task      Task `gorm:"foreignKey:TaskID"`
	User      User `gorm:"foreignKey:UserID"`
}
