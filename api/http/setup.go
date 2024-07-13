package http

import (
	"context"
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/middlerwares"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/routes"
	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/GoBootCamp-Group1/Task-Management/config"
	_ "github.com/GoBootCamp-Group1/Task-Management/docs"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func Run(cfg config.Server, app *app.Container) {
	fiberApp := fiber.New()
	fiberApp.Get("/swagger/*", swagger.HandlerDefault)

	api := fiberApp.Group("/api/v1", middlerwares.SetUserContext())

	//add basic role to db
	if err := app.RoleService().InitRolesInDb(context.Background()); err != nil {
		log.ErrorLog.Fatal(err)
	}

	// register global routes
	routes.InitAuthRoutes(&api, app)
	routes.InitBoardRoutes(&api, app, cfg)
	routes.InitTaskRoutes(&api, app, cfg)
	routes.InitColumnRoutes(&api, app, cfg)
	routes.InitNotificationRoutes(&api, app, cfg)
	routes.InitRoleRoutes(&api, app, cfg)

	// run server
	err := fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.HttpPort))
	if err != nil {
		log.ErrorLog.Fatal(err)
	}
	log.InfoLog.Println("Starting the application...")

}

func userRoleChecker() fiber.Handler {
	return middlerwares.RoleChecker("user")
}
