package routes

import (
	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers"
	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/gofiber/fiber/v2"
)

func InitAuthRoutes(router *fiber.Router, app *app.Container) {
	(*router).Post("/signup", handlers.SignUpUser(app.UserService()))
	(*router).Post("/login", handlers.LoginUser(app.AuthService()))
	(*router).Get("/refresh", handlers.RefreshCreds(app.AuthService()))
}
