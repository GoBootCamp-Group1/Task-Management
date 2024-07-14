package routes

import (
	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/middlerwares"
	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/GoBootCamp-Group1/Task-Management/config"
	"github.com/gofiber/fiber/v2"
)

func InitColumnRoutes(router *fiber.Router, container *app.Container, cfg config.Server) {
	columnGroup := (*router).Group("/columns", middlerwares.Auth([]byte(cfg.TokenSecret)))

	columnGroup.Post("/:boardId", handlers.CreateColumn(container.ColumnService()))
	columnGroup.Get("/:id", handlers.GetColumnByID(container.ColumnService()))
	columnGroup.Get("/:boardId/all", handlers.GetAllColumns(container.ColumnService()))
	columnGroup.Put("/:id", handlers.UpdateColumn(container.ColumnService()))
	columnGroup.Put("/:id/move", handlers.MoveColumn(container.ColumnService()))
	columnGroup.Put("/:id/final", handlers.ChangeFinalColumn(container.ColumnService()))
	columnGroup.Delete("/:id", handlers.DeleteColumn(container.ColumnService()))

}
