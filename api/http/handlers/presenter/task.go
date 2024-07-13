package presenter

import (
	"encoding/json"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/google/uuid"
	"time"
)

type TaskPresenter struct {
	ID            uint                     `json:"id"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
	OrderPosition int                      `json:"order_position"`
	Name          string                   `json:"name"`
	Description   string                   `json:"description"`
	StartDateTime *time.Time               `json:"start_datetime"`
	EndDateTime   *time.Time               `json:"end_datetime"`
	StoryPoint    int                      `json:"story_point"`
	Additional    json.RawMessage          `json:"additional"`
	Creator       *UserPresenter           `json:"creator"`
	Column        *ColumnOutBoundPresenter `json:"column"`
	Parent        *TaskPresenter           `json:"parent"`
	Assignee      *UserPresenter           `json:"assignee"`
}

func NewTaskPresenter(task *domains.Task) *TaskPresenter {

	var creator *UserPresenter
	if task.Creator != nil {
		creator = NewUserPresenter(task.Creator)
	}

	var column *ColumnOutBoundPresenter
	if task.Column != nil {
		column = NewColumnOutBoundPresenter(task.Column)
	}

	//var parent *TaskPresenter
	//if task.Parent != nil {
	//	parent = NewTaskPresenter(task.Parent)
	//}

	var assignee *UserPresenter
	if task.Assignee != nil {
		assignee = NewUserPresenter(task.Assignee)
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
		Column:        column,
		//Parent:        parent,
		Assignee: assignee,
	}
}

type TaskChildPresenter struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ColumnName    string    `json:"column_name"`
	ColumnIsFinal bool      `json:"column_is_final"`
}

func NewTaskChildPresenter(task *domains.TaskChild) *TaskChildPresenter {
	return &TaskChildPresenter{
		ID:            task.ID,
		CreatedAt:     task.CreatedAt,
		UpdatedAt:     task.UpdatedAt,
		Name:          task.Name,
		Description:   task.Description,
		ColumnName:    task.ColumnName,
		ColumnIsFinal: task.ColumnIsFinal,
	}
}

type TaskCommentPresenter struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Comment   string         `json:"comment"`
	User      *UserPresenter `json:"user"`
}

func NewTaskCommentPresenter(comment *domains.TaskComment) *TaskCommentPresenter {

	var user *UserPresenter
	if comment.User != nil {
		user = NewUserPresenter(comment.User)
	}

	return &TaskCommentPresenter{
		ID:        comment.ID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		Comment:   comment.Comment,
		User:      user,
	}
}
