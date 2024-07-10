package routes

import (
	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/middlerwares"
	"github.com/GoBootCamp-Group1/Task-Management/config"

	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/gofiber/fiber/v2"
)

func InitBoardRoutes(router *fiber.Router, app *app.Container, cfg config.Server) {

	boardGroup := (*router).Group("/boards", middlerwares.Auth([]byte(cfg.TokenSecret)))

	boardGroup.Post("", handlers.CreateBoard(app.BoardService()))
	boardGroup.Put("/:id", handlers.UpdateBoard(app.BoardService()))
	boardGroup.Get("/:id", handlers.GetBoardByID(app.BoardService()))
	boardGroup.Delete("/:id", handlers.DeleteBoard(app.BoardService()))
}
