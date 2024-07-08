package route

import (
	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers"
	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/gofiber/fiber/v2"
)

func InitBoardRoutes(router *fiber.Router, app *app.Container) {
	(*router).Post("/boards", handlers.CreateBoard(app.BoardService()))
	(*router).Put("/boards/:id", handlers.UpdateBoard(app.BoardService()))
	(*router).Get("/boards/:id", handlers.GetBoardByID(app.BoardService()))
	(*router).Delete("/boards/:id", handlers.DeleteBoard(app.BoardService()))
}
