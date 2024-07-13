package routes

import (
	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/middlerwares"
	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/GoBootCamp-Group1/Task-Management/config"
	"github.com/gofiber/fiber/v2"
)

func InitRoleRoutes(router *fiber.Router, container *app.Container, cfg config.Server) {
	roleGroup := (*router).Group("/roles", middlerwares.Auth([]byte(cfg.TokenSecret)))
	roleGroup.Post("", handlers.CreateRole(container.RoleService()))
	roleGroup.Put("/:id", handlers.UpdateRole(container.RoleService()))
	roleGroup.Get("/:id", handlers.GetRoleByID(container.RoleService()))
	roleGroup.Delete("/:id", handlers.DeleteRole(container.RoleService()))
}
