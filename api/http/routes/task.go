package routes

import (
	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/middlerwares"
	"github.com/GoBootCamp-Group1/Task-Management/config"

	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/gofiber/fiber/v2"
)

func InitTaskRoutes(router *fiber.Router, app *app.Container, cfg config.Server) {

	taskGroup := (*router).Group("/boards/:boardID/tasks", middlerwares.Auth([]byte(cfg.TokenSecret)))

	taskGroup.Post("/", handlers.CreateTask(app.TaskService()))
	taskGroup.Put("/:id", handlers.UpdateTask(app.TaskService()))
	taskGroup.Get("/", handlers.GetTasksByBoardID(app.TaskService()))
	taskGroup.Get("/:id", handlers.GetTaskByID(app.TaskService()))
	taskGroup.Get("/:id/children", handlers.GetTaskChildren(app.TaskService()))
	taskGroup.Delete("/:id", handlers.DeleteTask(app.TaskService()))

	taskGroup.Patch("/:id/column", handlers.ChangeTaskColumn(app.TaskService()))

	taskGroup.Post("/:taskID/dependencies/:dependentTaskID", handlers.AddTaskDependency(app.TaskService()))
	taskGroup.Delete("/:taskID/dependencies/:dependentTaskID", handlers.RemoveTaskDependency(app.TaskService()))
	taskGroup.Get("/:taskID/dependencies", handlers.GetTaskDependencies(app.TaskService()))

	taskGroup.Post("/:taskID/comments", handlers.CreateTaskComment(app.TaskService()))
	taskGroup.Get("/:taskID/comments", handlers.TaskCommentsList(app.TaskService()))
	taskGroup.Get("/:taskID/comments/:id", handlers.GetCommentByID(app.TaskService()))
	taskGroup.Delete("/:taskID/comments/:id", handlers.DeleteComment(app.TaskService()))
}
