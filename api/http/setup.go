package http

import (
	"fmt"
	"log"

	"github.com/GoBootCamp-Group1/Task-Management/api/http/middlerwares"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/routes"
	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/GoBootCamp-Group1/Task-Management/config"
	_ "github.com/GoBootCamp-Group1/Task-Management/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func Run(cfg config.Server, app *app.Container) {
	fiberApp := fiber.New()
	fiberApp.Get("/swagger/*", swagger.HandlerDefault)

	api := fiberApp.Group("/api/v1", middlerwares.SetUserContext())

	// register global routes
	routes.InitAuthRoutes(&api, app)
	routes.InitBoardRoutes(&api, app, cfg)
	routes.InitTaskRoutes(&api, app, cfg)
	routes.InitColumnRoutes(&api, app, cfg)

	// run server
	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.HttpPort)))
}

func userRoleChecker() fiber.Handler {
	return middlerwares.RoleChecker("user")
}
