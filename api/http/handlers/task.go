package handlers

import (
	"encoding/json"
	"errors"
	"time"

	"fmt"

	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers/presenter"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/services"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/log"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type TaskRequest struct {
	Name          string          `json:"name" validate:"required,min=3,max=50" example:"new task"`
	ColumnID      uint            `json:"column_id" validate:"required,gte=1" example:"1"`
	ParentID      *uint           `json:"parent_id,omitempty" validate:"omitempty,gte=1" example:"1"`
	AssigneeID    *uint           `json:"assignee_id,omitempty" validate:"omitempty,gte=1" example:"1"`
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
	ErrTaskNotFound               = fiber.NewError(fiber.StatusNotFound, "Task not found")
	ErrNoTaskFound                = fiber.NewError(fiber.StatusNotFound, "No tasks found")
)

// CreateTask creates a new task
// @Summary Create Task
// @Description creates a task
// @Tags Task
// @Accept  json
// @Produce json
// @Param   body  body      TaskRequest  true  "Create Task"
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
			log.ErrorLog.Printf("Error loading user: %v\n", errUserID)
			return OldSendError(c, errUserID, fiber.StatusInternalServerError)
		}

		//Get Board ID
		boardID, errBoardID := c.ParamsInt("boardID")
		if errBoardID != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", errBoardID)
			return OldSendError(c, ErrInvalidBoardIDParam, fiber.StatusBadRequest)
		}

		//datetime parse
		startDateTime, errInvalidStartDt := time.Parse(dateTimeLayout, input.StartDateTime)
		if errInvalidStartDt != nil {
			log.ErrorLog.Printf("Error invalid time: %v\n", errInvalidStartDt)
			return OldSendError(c, ErrInvalidStartDatetimeLayout, fiber.StatusBadRequest)
		}
		endDateTime, errInvalidEndDt := time.Parse(dateTimeLayout, input.EndDateTime)
		if errInvalidEndDt != nil {
			log.ErrorLog.Printf("Error invalid time: %v\n", errInvalidEndDt)
			return OldSendError(c, ErrInvalidEndDatetimeLayout, fiber.StatusBadRequest)
		}

		taskModel := domains.Task{
			CreatedBy:     userID,
			BoardID:       uint(boardID),
			ParentID:      input.ParentID,
			AssigneeID:    input.AssigneeID,
			ColumnID:      input.ColumnID,
			OrderPosition: input.OrderPosition,
			Name:          input.Name,
			Description:   input.Description,
			StartDateTime: &startDateTime,
			EndDateTime:   &endDateTime,
			StoryPoint:    input.StoryPoint,
			Additional:    input.Additional,
		}

		createdTask, err := taskService.CreateTask(c.Context(), &taskModel)
		if err != nil {
			log.ErrorLog.Printf("Error creating task: %v\n", err)
			return OldSendError(c, err, fiber.StatusInternalServerError)
		}
		log.InfoLog.Println("Task created successfully")

		return SendSuccessResponse(
			c,
			"Successfully created.",
			presenter.NewTaskPresenter(createdTask),
		)
	}
}

// GetTaskByID get a task
// @Summary Get Task
// @Description gets a task
// @Tags Task
// @Produce json
// @Param   id      path     string  true  "Task ID"
// @Success 200 {object} Response
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /boards/{boardID}/tasks/{taskID} [get]
// @Security ApiKeyAuth
func GetTaskByID(taskService *services.TaskService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			log.ErrorLog.Printf("Error parsing task id: %v\n", err)
			return OldSendError(c, err, fiber.StatusBadRequest)
		}

		//Get User ID
		userID, errUserID := utils.GetUserID(c)
		if errUserID != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", errUserID)
			return OldSendError(c, errUserID, fiber.StatusInternalServerError)
		}

		boardID, err := c.ParamsInt("boardID")
		if err != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", err)
			return OldSendError(c, err, fiber.StatusBadRequest)
		}

		task, err := taskService.GetTaskByID(c.Context(), userID, uint(boardID), uint(id))
		if err != nil {
			log.ErrorLog.Printf("Error getting task: %v\n", err)
			return OldSendError(c, err, fiber.StatusInternalServerError)
		}

		if task == nil {
			log.ErrorLog.Printf("Error getting task: %v\n", ErrTaskNotFound)
			return OldSendError(c, ErrTaskNotFound, fiber.StatusNotFound)
		}
		log.InfoLog.Println("Task loaded successfully")

		return SendSuccessResponse(
			c,
			"Successfully fetched.",
			presenter.NewTaskPresenter(task),
		)
	}
}

// GetTasksByBoardID get tasks for a board
// @Summary Get Tasks
// @Description gets tasks for a board
// @Tags Task
// @Produce json
// @Param   boardID  path     string  true  "Board ID"
// @Success 200 {array} Response
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /boards/{boardID}/tasks [get]
// @Security ApiKeyAuth
func GetTasksByBoardID(taskService *services.TaskService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		boardID, err := c.ParamsInt("boardID")
		if err != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", err)
			return OldSendError(c, err, fiber.StatusBadRequest)
		}

		//Get User ID
		userID, errUserID := utils.GetUserID(c)
		if errUserID != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", errUserID)
			return OldSendError(c, errUserID, fiber.StatusInternalServerError)
		}

		// init variables for pagination
		page, pageSize := PageAndPageSize(c)

		tasks, total, err := taskService.GetTasksByBoardID(c.Context(), userID, uint(boardID), uint(page), uint(pageSize))
		if err != nil {
			log.ErrorLog.Printf("Error gettings tasks: %v\n", err)
			return OldSendError(c, err, fiber.StatusInternalServerError)
		}

		if len(tasks) == 0 {
			log.ErrorLog.Printf("Error getting tasks: %v\n", ErrNoTaskFound)
			return OldSendError(c, ErrNoTaskFound, fiber.StatusNotFound)
		}

		//generate response data
		taskPresenters := make([]*presenter.TaskPresenter, len(tasks))
		for i, task := range tasks {
			taskPresenters[i] = presenter.NewTaskPresenter(&task)
		}
		log.InfoLog.Println("Tasks loaded successfully")

		return SendSuccessPaginateResponse(
			c,
			"Successfully fetched.",
			taskPresenters,
			uint(page),
			uint(pageSize),
			total,
		)
	}
}

// UpdateTask updates a new task
// @Summary Update Task
// @Description updates a task
// @Tags Task
// @Accept  json
// @Produce json
// @Param   id      path     string  true  "Task ID"
// @Param   body  body      TaskRequest  true  "Update Task"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /boards/{boardID}/tasks/{taskID} [put]
// @Security ApiKeyAuth
func UpdateTask(taskService *services.TaskService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var input TaskRequest

		err := ValidateAndFill(c, &input)
		if err != nil {
			return err
		}

		id, err := c.ParamsInt("id")
		if err != nil {
			log.ErrorLog.Printf("Error parsing task id: %v\n", err)
			return OldSendError(c, err, fiber.StatusBadRequest)
		}

		boardID, err := c.ParamsInt("boardID")
		if err != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", err)
			return OldSendError(c, err, fiber.StatusBadRequest)
		}

		//Get User ID
		userID, errUserID := utils.GetUserID(c)
		if errUserID != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", errUserID)
			return OldSendError(c, errUserID, fiber.StatusInternalServerError)
		}

		//datetime parse
		startDateTime, errInvalidStartDt := time.Parse(dateTimeLayout, input.StartDateTime)
		if errInvalidStartDt != nil {
			log.ErrorLog.Printf("Error invalid time: %v\n", errInvalidStartDt)
			return OldSendError(c, ErrInvalidStartDatetimeLayout, fiber.StatusBadRequest)
		}
		endDateTime, errInvalidEndDt := time.Parse(dateTimeLayout, input.EndDateTime)
		if errInvalidEndDt != nil {
			log.ErrorLog.Printf("Error invalid time: %v\n", errInvalidStartDt)
			return OldSendError(c, ErrInvalidEndDatetimeLayout, fiber.StatusBadRequest)
		}

		taskModel := domains.Task{
			ID:            uint(id),
			ParentID:      input.ParentID,
			AssigneeID:    input.AssigneeID,
			ColumnID:      input.ColumnID,
			OrderPosition: input.OrderPosition,
			Name:          input.Name,
			Description:   input.Description,
			StartDateTime: &startDateTime,
			EndDateTime:   &endDateTime,
			StoryPoint:    input.StoryPoint,
			Additional:    input.Additional,
		}

		updatedTask, err := taskService.UpdateTask(c.Context(), userID, uint(boardID), &taskModel)
		if err != nil {
			log.ErrorLog.Printf("Error updating task: %v\n", err)
			return OldSendError(c, err, fiber.StatusInternalServerError)
		}

		if updatedTask == nil {
			return OldSendError(c, ErrTaskNotFound, fiber.StatusNotFound)
		}

		log.InfoLog.Println("Task updated successfully")

		return SendSuccessResponse(
			c,
			"Successfully updated.",
			presenter.NewTaskPresenter(updatedTask),
		)
	}
}

// DeleteTask delete a task
// @Summary Delete Task
// @Description deleted a task
// @Tags Task
// @Produce json
// @Param   id      path     string  true  "Task ID"
// @Success 204
// @Failure 400
// @Failure 500
// @Router /boards/{boardID}/tasks/{id} [delete]
// @Security ApiKeyAuth
func DeleteTask(taskService *services.TaskService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			log.ErrorLog.Printf("Error parsing task id: %v\n", err)
			return OldSendError(c, err, fiber.StatusBadRequest)
		}

		//Get User ID
		userID, errUserID := utils.GetUserID(c)
		if errUserID != nil {
			log.ErrorLog.Printf("Error loading task: %v\n", errUserID)
			return OldSendError(c, errUserID, fiber.StatusInternalServerError)
		}

		err = taskService.DeleteTask(c.Context(), userID, uint(id))
		if err != nil {
			log.ErrorLog.Printf("Error deleting task: %v\n", err)
			return OldSendError(c, err, fiber.StatusInternalServerError)
		}
		log.InfoLog.Println("Task deleted successfully")

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// GetTaskChildren get a list of task children
// @Summary Get TaskChildren
// @Description get list of a task children
// @Tags Task
// @Produce json
// @Param   id      path     string  true  "Task ID"
// @Success 204
// @Failure 400
// @Failure 500
// @Router /boards/{boardID}/tasks/{id}/children [get]
// @Security ApiKeyAuth
func GetTaskChildren(taskService *services.TaskService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			log.ErrorLog.Printf("Error parsing task id: %v\n", err)
			return OldSendError(c, err, fiber.StatusBadRequest)
		}

		boardID, err := c.ParamsInt("boardID")
		if err != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", err)
			return OldSendError(c, err, fiber.StatusBadRequest)
		}

		//Get User ID
		userID, errUserID := utils.GetUserID(c)
		if errUserID != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", errUserID)
			return OldSendError(c, errUserID, fiber.StatusInternalServerError)
		}

		children, errFetchChildren := taskService.GetTaskChildren(c.Context(), userID, uint(boardID), uint(id))
		if errFetchChildren != nil {
			log.ErrorLog.Printf("Error loading children: %v\n", errFetchChildren)
			return OldSendError(c, errFetchChildren, fiber.StatusInternalServerError)
		}

		//generate response data
		childrenPresenters := make([]*presenter.TaskChildPresenter, len(children))
		for i, task := range children {
			childrenPresenters[i] = presenter.NewTaskChildPresenter(&task)
		}

		return SendSuccessResponse(
			c,
			"Successfully fetched.",
			childrenPresenters,
		)
	}
}

type ColumnChangeRequest struct {
	NewColumnID uint `json:"new_column_id" validate:"required,gte=1" example:"1"`
}

// ChangeTaskColumn updates task column
// @Summary Update Task column
// @Description updates a task column
// @Tags Task
// @Accept  json
// @Produce json
// @Param   boardID      path     string  true  "Board ID"
// @Param   taskID      path     string  true  "Task ID"
// @Param   body  body      ColumnChangeRequest  true  "Update Task Column"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /boards/{boardID}/tasks/{taskID}/column [patch]
// @Security ApiKeyAuth
func ChangeTaskColumn(taskService *services.TaskService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var input ColumnChangeRequest

		err := ValidateAndFill(c, &input)
		if err != nil {
			return err
		}

		id, err := c.ParamsInt("id")
		if err != nil {
			log.ErrorLog.Printf("Error parsing task id: %v\n", err)
			return OldSendError(c, err, fiber.StatusBadRequest)
		}

		boardID, err := c.ParamsInt("boardID")
		if err != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", err)
			return OldSendError(c, err, fiber.StatusBadRequest)
		}

		//Get User ID
		userID, errUserID := utils.GetUserID(c)
		if errUserID != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", errUserID)
			return OldSendError(c, errUserID, fiber.StatusInternalServerError)
		}

		taskModel := domains.Task{
			ID:      uint(id),
			BoardID: uint(boardID),
		}

		updatedTask, err := taskService.ChangeTaskColumn(c.Context(), userID, &taskModel, input.NewColumnID)
		if err != nil {
			log.ErrorLog.Printf("Error updating task: %v\n", err)
			return OldSendError(c, err, fiber.StatusInternalServerError)
		}

		if updatedTask == nil {
			return OldSendError(c, ErrTaskNotFound, fiber.StatusNotFound)
		}

		log.InfoLog.Println("Task updated successfully")

		return SendSuccessResponse(
			c,
			"Successfully updated.",
			presenter.NewTaskPresenter(updatedTask),
		)
	}
}

// AddTaskDependency adds a dependency between two tasks
// @Summary Add Task Dependency
// @Description adds a dependency between two tasks
// @Tags Task
// @Accept  json
// @Param   taskID  path   uint  true  "Task ID"
// @Param   dependentTaskID  path   uint  true  "Dependent Task ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /boards/{boardID}/tasks/{taskID}/dependencies/{dependentTaskID} [post]
// @Security ApiKeyAuth
func AddTaskDependency(taskService *services.TaskService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse taskID and dependentTaskID from path parameters
		taskID, err := c.ParamsInt("taskID")
		if err != nil {
			log.ErrorLog.Printf("Error parsing taskID: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		dependentTaskID, err := c.ParamsInt("dependentTaskID")
		if err != nil {
			log.ErrorLog.Printf("Error parsing dependentTaskID: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		// Add task dependency
		err = taskService.AddTaskDependency(c.Context(), uint(taskID), uint(dependentTaskID))
		if err != nil {
			log.ErrorLog.Printf("Error adding task dependency: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		log.InfoLog.Printf("Added dependency from task #%d to task #%d", taskID, dependentTaskID)
		return SendSuccessResponse(c, "Task dependency added successfully", nil)
	}
}

// RemoveTaskDependency removes a dependency between two tasks
// @Summary Remove Task Dependency
// @Description removes a dependency between two tasks
// @Tags Task
// @Accept  json
// @Param   taskID  path   uint  true  "Task ID"
// @Param   dependentTaskID  path   uint  true  "Dependent Task ID"
// @Success 204
// @Failure 400
// @Failure 500
// @Router /boards/{boardID}/tasks/{taskID}/dependencies/{dependentTaskID} [delete]
// @Security ApiKeyAuth
func RemoveTaskDependency(taskService *services.TaskService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse taskID and dependentTaskID from path parameters
		taskID, err := c.ParamsInt("taskID")
		if err != nil {
			log.ErrorLog.Printf("Error parsing taskID: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		dependentTaskID, err := c.ParamsInt("dependentTaskID")
		if err != nil {
			log.ErrorLog.Printf("Error parsing dependentTaskID: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		// Remove task dependency
		err = taskService.RemoveTaskDependency(c.Context(), uint(taskID), uint(dependentTaskID))
		if err != nil {
			log.ErrorLog.Printf("Error removing task dependency: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		log.InfoLog.Printf("Removed dependency from task #%d to task #%d", taskID, dependentTaskID)
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// GetTaskDependencies retrieves the dependencies for a given task
// @Summary Get Task Dependencies
// @Description retrieves the dependencies for a given task
// @Tags Task
// @Accept  json
// @Param   taskID  path   uint  true  "Task ID"
// @Success 200 {array} domains.TaskDependency
// @Failure 400
// @Failure 500
// @Router /boards/{boardID}/tasks/{taskID}/dependencies [get]
// @Security ApiKeyAuth
func GetTaskDependencies(taskService *services.TaskService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse taskID from path parameters
		taskID, err := c.ParamsInt("taskID")
		if err != nil {
			log.ErrorLog.Printf("Error parsing taskID: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		// Get task dependencies
		taskDependencies, err := taskService.GetTaskDependencies(c.Context(), uint(taskID))
		if err != nil {
			log.ErrorLog.Printf("Error retrieving task dependencies: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		msg := fmt.Sprintf("Retrieved dependencies for task #%d", taskID)
		log.InfoLog.Println(msg)
		return SendSuccessResponse(c, msg, taskDependencies)
	}
}
