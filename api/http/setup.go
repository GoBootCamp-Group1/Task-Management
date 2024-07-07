package http

import (
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/middlerwares"
	"github.com/GoBootCamp-Group1/Task-Management/config"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/service"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"log"
)

func Run(cfg config.Server, app *service.AppContainer) {
	fiberApp := fiber.New()
	api := fiberApp.Group("/api/v1", middlerwares.SetUserContext())

	// register global routes
	registerGlobalRoutes(api, app)

	secret := []byte(cfg.TokenSecret)

	// registering users APIs
	registerUsersAPI(api, app.UserService(), secret)

	// run server
	log.Fatal(fiberApp.Listen(fmt.Sprintf("%s:%d", cfg.Host, cfg.HttpPort)))
}

func registerUsersAPI(router fiber.Router, _ *service.UserService, secret []byte) {
	userGroup := router.Group("/users", middlerwares.Auth(secret), middlerwares.RoleChecker("user", "admin"))

	userGroup.Get("/folan", func(c *fiber.Ctx) error {
		claims := c.Locals(jwt.UserClaimKey).(*jwt.UserClaims)

		return c.JSON(map[string]any{
			"user_id": claims.UserID,
			"role":    claims.Role,
		})
	})
}

func registerGlobalRoutes(router fiber.Router, app *service.AppContainer) {
	router.Post("/signup", handlers.SignUpUser(app.UserService()))
	router.Post("/login", handlers.LoginUser(app.AuthService()))
	router.Get("/refresh", handlers.RefreshCreds(app.AuthService()))
}

func userRoleChecker() fiber.Handler {
	return middlerwares.RoleChecker("user")
}
