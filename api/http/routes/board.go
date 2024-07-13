package routes

import (
	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/middlerwares"
	"github.com/GoBootCamp-Group1/Task-Management/config"

	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/gofiber/fiber/v2"
)

func InitBoardRoutes(router *fiber.Router, cotainer *app.Container, cfg config.Server) {

	boardGroup := (*router).Group("/boards", middlerwares.Auth([]byte(cfg.TokenSecret)))

	boardGroup.Post("", handlers.CreateBoard(cotainer.BoardService()))
	boardGroup.Put("/:id", handlers.UpdateBoard(cotainer.BoardService()))
	boardGroup.Get("/:id", handlers.GetBoardByID(cotainer.BoardService()))
	boardGroup.Delete("/:id", handlers.DeleteBoard(cotainer.BoardService()))

	boardGroup.Post("/:id/add-user", handlers.InviteUserToBoard(cotainer.BoardService()))
	boardGroup.Delete("/:board_id/users/:user_id", handlers.RemoveUserFromBoard(cotainer.BoardService()))
	boardGroup.Put("/:board_id/users/:user_id", handlers.ChangeUserRoleInBoard(cotainer.BoardService()))
}
