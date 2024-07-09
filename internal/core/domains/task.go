package domains

import (
	"encoding/json"
	"time"
)

type Task struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	CreatedBy uint

	BoardID       uint
	ParentID      *uint
	AssigneeID    *uint
	ColumnID      uint
	OrderPosition int
	Name          string
	Description   string
	StartDateTime *time.Time
	EndDateTime   *time.Time
	StoryPoint    *int
	Additional    json.RawMessage

	//Board    Board
	//Column   Column
	//Parent   Task
	//Assignee User
	//Creator  User
}
