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
	ParentID      uint
	AssigneeID    uint
	ColumnID      uint
	OrderPosition int
	Name          string
	Description   string
	StartDateTime time.Time
	EndDateTime   time.Time
	StoryPoint    int
	Additional    json.RawMessage

	//Board    Board
	//Column   Column
	//Parent   Task
	//Assignee User
	//Creator  User
}
