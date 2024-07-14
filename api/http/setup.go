package http

import (
	"context"
	"fmt"
	"time"

	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/middlerwares"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/routes"
	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/GoBootCamp-Group1/Task-Management/config"
	_ "github.com/GoBootCamp-Group1/Task-Management/docs"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
)

func Run(cfg config.Server, app *app.Container) {
	fiberApp := fiber.New()
	fiberApp.Get("/swagger/*", swagger.HandlerDefault)

	rateLimit(cfg, fiberApp)
	corsLimit(cfg, fiberApp)

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

func rateLimit(cfg config.Server, fiberApp *fiber.App) fiber.Router {
	return fiberApp.Use(limiter.New(limiter.Config{
		Max:        cfg.MaxRateLimit,
		Expiration: time.Duration(cfg.RateLimitExpiration) * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			log.ErrorLog.Println("User reached request limit!")
			return handlers.SendError(c, fiber.NewError(fiber.StatusTooManyRequests, "Too many requests"))
		},
	}))
}

func corsLimit(cfg config.Server, fiberApp *fiber.App) fiber.Router {
	return fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins,
		AllowMethods: cfg.AllowedMethods,
		AllowHeaders: cfg.AllowedMHeaders,
	}))
}

func userRoleChecker() fiber.Handler {
	return middlerwares.RoleChecker("user")
}
