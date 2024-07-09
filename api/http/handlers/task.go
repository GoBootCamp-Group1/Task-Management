package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/services"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/utils"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

type TaskRequest struct {
	Name          string          `json:"name" validate:"required,min=3,max=50" example:"new task"`
	ColumnID      uint            `json:"column_id" validate:"required,gte=1" example:"1"`
	ParentID      uint            `json:"parent_id,omitempty" validate:"gte=1" example:"1"`
	AssigneeID    uint            `json:"assignee_id,omitempty" validate:"gte=1" example:"1"`
	OrderPosition int             `json:"order_position" validate:"required,number" example:"1"`
	Description   string          `json:"description" validate:"required,min=1,max=2000" example:"This is the description"`
	StartDateTime string          `json:"start_datetime" validate:"required" example:"2020-01-01 16:30:00"`
	EndDateTime   string          `json:"end_datetime" validate:"required" example:"2020-01-01 16:30:00"`
	StoryPoint    int             `json:"story_point" validate:"required,number" example:"1"`
	Additional    json.RawMessage `json:"additional,omitempty" validate:"json"`
}

var dateTimeLayout = "2006-01-02 15:04:05"

var (
	ErrInvalidBoardIDParam        = errors.New("invalid board id")
	ErrInvalidStartDatetimeLayout = errors.New("invalid start datetime format, example: " + dateTimeLayout)
	ErrInvalidEndDatetimeLayout   = errors.New("invalid end datetime format, example: " + dateTimeLayout)
)

// CreateTask creates a new task
// @Summary Create Task
// @Description creates a task
// @Tags Create Task
// @Accept  json
// @Produce json
// @Param   body  body      CreateTaskRequest  true  "Create Task"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /tasks [post]
// @Security ApiKeyAuth
func CreateTask(taskService *services.TaskService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var input TaskRequest

		err := ValidateAndFill(c, &input)
		if err != nil {
			return err
		}

		//Get User ID
		userID, errUserID := utils.GetUserID(c)
		if errUserID != nil {
			return SendError(c, errUserID, fiber.StatusInternalServerError)
		}

		//Get Board ID
		boardID, errBoardID := c.ParamsInt("boardID")
		if errBoardID != nil {
			return SendError(c, ErrInvalidBoardIDParam, fiber.StatusBadRequest)
		}

		//datetime parse
		startDateTime, errInvalidStartDt := time.Parse(dateTimeLayout, input.StartDateTime)
		if errInvalidStartDt != nil {
			return SendError(c, ErrInvalidStartDatetimeLayout, fiber.StatusBadRequest)
		}
		endDateTime, errInvalidEndDt := time.Parse(dateTimeLayout, input.EndDateTime)
		if errInvalidEndDt != nil {
			return SendError(c, ErrInvalidEndDatetimeLayout, fiber.StatusBadRequest)
		}

		taskModel := domains.Task{
			CreatedBy:     userID,
			BoardID:       uint(boardID),
			ParentID:      &input.ParentID,
			AssigneeID:    &input.AssigneeID,
			ColumnID:      input.ColumnID,
			OrderPosition: input.OrderPosition,
			Name:          input.Name,
			Description:   input.Description,
			StartDateTime: &startDateTime,
			EndDateTime:   &endDateTime,
			StoryPoint:    &input.StoryPoint,
			Additional:    input.Additional,
		}

		err = taskService.CreateTask(c.Context(), &taskModel)
		if err != nil {
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).JSON(map[string]any{
			"data":    taskModel,
			"message": "Successfully created.",
		})
	}
}

// GetTaskByID get a task
// @Summary Get Task
// @Description gets a task
// @Tags Get Task
// @Produce json
// @Param   id      path     string  true  "Task ID"
// @Success 200 {object} domains.Task
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /tasks/{id} [get]
// @Security ApiKeyAuth
func GetTaskByID(taskService *services.TaskService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		task, err := taskService.GetTaskByID(c.Context(), uint(id))
		if err != nil {
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		if task == nil {
			return SendError(c, fiber.NewError(fiber.StatusNotFound, "Task not found"), fiber.StatusNotFound)
		}

		return c.JSON(task)
	}
}

type UpdateTaskRequest struct {
	Name      string `json:"name" validate:"required,min=3,max=50,excludesall=;" example:"new task"`
	IsPrivate bool   `json:"is_private" example:"false"`
}

// UpdateTask updates a new task
// @Summary Update Task
// @Description updates a task
// @Tags Update Task
// @Accept  json
// @Produce json
// @Param   id      path     string  true  "Task ID"
// @Param   body  body      UpdateTaskRequest  true  "Update Task"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /tasks/{id} [put]
// @Security ApiKeyAuth
func UpdateTask(taskService *services.TaskService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		validate := validation.NewValidator()
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}
		var input UpdateTaskRequest

		if err = c.BodyParser(&input); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err = validate.Struct(input)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		taskModel := domains.Task{
			ID:   uint(id),
			Name: input.Name,
		}

		err = taskService.UpdateTask(c.Context(), &taskModel)
		if err != nil {
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		return SendSuccessResponse(c, "task")
	}
}

// DeleteTask delete a task
// @Summary Delete Task
// @Description deleted a task
// @Tags Delete Task
// @Produce json
// @Param   id      path     string  true  "Task ID"
// @Success 204
// @Failure 400
// @Failure 500
// @Router /tasks/{id} [delete]
// @Security ApiKeyAuth
func DeleteTask(taskService *services.TaskService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err = taskService.DeleteTask(c.Context(), uint(id))
		if err != nil {
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
