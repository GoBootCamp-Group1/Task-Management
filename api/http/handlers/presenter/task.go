package presenter

import (
	"encoding/json"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"time"
)

type TaskPresenter struct {
	ID            uint            `json:"id"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	OrderPosition int             `json:"order_position"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	StartDateTime *time.Time      `json:"start_datetime"`
	EndDateTime   *time.Time      `json:"end_datetime"`
	StoryPoint    int             `json:"story_point"`
	Additional    json.RawMessage `json:"additional"`
	Creator       *UserPresenter  `json:"creator"`
	Parent        *TaskPresenter  `json:"parent"`
}

func NewTaskPresenter(task *domains.Task) *TaskPresenter {

	var creator *UserPresenter
	if task.Creator != nil {
		creator = NewUserPresenter(task.Creator)
	}

	var parent *TaskPresenter
	if task.Parent != nil {
		parent = NewTaskPresenter(task.Parent)
	}

	return &TaskPresenter{
		ID:            task.ID,
		CreatedAt:     task.CreatedAt,
		UpdatedAt:     task.UpdatedAt,
		OrderPosition: task.OrderPosition,
		Name:          task.Name,
		Description:   task.Description,
		StartDateTime: task.StartDateTime,
		EndDateTime:   task.EndDateTime,
		StoryPoint:    task.StoryPoint,
		Additional:    task.Additional,
		Creator:       creator,
		Parent:        parent,
	}
}
