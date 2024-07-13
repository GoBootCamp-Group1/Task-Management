package routes

import (
	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/middlerwares"
	"github.com/GoBootCamp-Group1/Task-Management/config"

	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/gofiber/fiber/v2"
)

func InitBoardRoutes(router *fiber.Router, container *app.Container, cfg config.Server) {

	boardGroup := (*router).Group("/boards", middlerwares.Auth([]byte(cfg.TokenSecret)))

	boardGroup.Post("", handlers.CreateBoard(container.BoardService()))
	boardGroup.Put("/:id", handlers.UpdateBoard(container.BoardService()))
	boardGroup.Get("/:id", handlers.GetBoardByID(container.BoardService()))
	boardGroup.Delete("/:id", handlers.DeleteBoard(container.BoardService()))

	boardGroup.Post("/:id/add-user", handlers.InviteUserToBoard(container.BoardService()))
	boardGroup.Delete("/:board_id/users/:user_id", handlers.RemoveUserFromBoard(container.BoardService()))
	boardGroup.Put("/:board_id/users/:user_id", handlers.ChangeUserRoleInBoard(container.BoardService()))
}
