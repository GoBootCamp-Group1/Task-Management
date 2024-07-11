package http

import (
	"context"
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"

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

	//add basic role to db
	createRolesInDb(app)
}

func userRoleChecker() fiber.Handler {
	return middlerwares.RoleChecker("user")
}

func createRolesInDb(cotainer *app.Container) {

	maintainerRoleEnum, _ := domains.ParseRole("Maintainer")
	editorRoleEnum, _ := domains.ParseRole("Editor")
	ownerRoleEnum, _ := domains.ParseRole("Owner")
	viewerRoleEnum, _ := domains.ParseRole("Viewer")
	maintainerRole := domains.Role{
		ID:          0,
		Name:        maintainerRoleEnum.String(),
		Description: "its maintainer role",
		Weight:      int(maintainerRoleEnum),
	}
	editorRole := domains.Role{
		ID:          0,
		Name:        editorRoleEnum.String(),
		Description: "its maintainer role",
		Weight:      int(editorRoleEnum),
	}
	ownerRole := domains.Role{
		ID:          0,
		Name:        ownerRoleEnum.String(),
		Description: "its maintainer role",
		Weight:      int(ownerRoleEnum),
	}
	viewerRole := domains.Role{
		ID:          0,
		Name:        viewerRoleEnum.String(),
		Description: "its maintainer role",
		Weight:      int(viewerRoleEnum),
	}
	//unhandled error
	_ = cotainer.RoleService().CreateRole(context.Background(), &maintainerRole)
	_ = cotainer.RoleService().CreateRole(context.Background(), &editorRole)
	_ = cotainer.RoleService().CreateRole(context.Background(), &ownerRole)
	_ = cotainer.RoleService().CreateRole(context.Background(), &viewerRole)
}
