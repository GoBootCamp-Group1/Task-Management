package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/middlerwares"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/routes"
	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/GoBootCamp-Group1/Task-Management/config"
	_ "github.com/GoBootCamp-Group1/Task-Management/docs"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
	"time"
)

func Run(cfg config.Server, app *app.Container) {
	fiberApp := fiber.New()
	fiberApp.Get("/swagger/*", swagger.HandlerDefault)

	fiberApp.Use(limiter.New(limiter.Config{
		Max:        cfg.MaxRateLimit,
		Expiration: time.Duration(cfg.RateLimitExpiration) * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			log.ErrorLog.Println("User reached request limit!")
			return handlers.SendError(c, errors.New("too many requests"), fiber.StatusTooManyRequests)
		},
	}))

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
	if err = app.RoleService().InitRolesInDb(context.Background()); err != nil {
		log.ErrorLog.Fatal(err)
	}
}

func userRoleChecker() fiber.Handler {
	return middlerwares.RoleChecker("user")
}
