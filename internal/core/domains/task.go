package domains

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID            uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	CreatedBy     uint
	BoardID       uint
	ParentID      *uint
	AssigneeID    *uint
	ColumnID      uint
	OrderPosition int
	Name          string
	Description   string
	StartDateTime *time.Time
	EndDateTime   *time.Time
	StoryPoint    int
	Additional    json.RawMessage
	Board         *Board
	Creator       *User
	Column        *Column
	//Parent        *Task
	Assignee *User
}

type TaskChild struct {
	ID            uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	ColumnID      uint
	OrderPosition int
	Name          string
	Description   string
	ColumnName    string
	ColumnIsFinal bool
}

type TaskComment struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	UserID    uint
	TaskID    uint
	Comment   string
	Task      *Task
	User      *User
}
