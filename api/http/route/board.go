package route

import (
	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers"
	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/gofiber/fiber/v2"
)

func InitBoardRoutes(router *fiber.Router, app *app.Container) {
	(*router).Post("/board", handlers.CreateBoard(app.BoardService()))
	(*router).Put("/board/:id", handlers.UpdateBoard(app.BoardService()))
	(*router).Get("/board/:id", handlers.GetBoardByID(app.BoardService()))
	(*router).Delete("/board/:id", handlers.DeleteBoard(app.BoardService()))
}
