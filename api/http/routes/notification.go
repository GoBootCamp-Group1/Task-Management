package routes

import (
	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/middlerwares"
	"github.com/GoBootCamp-Group1/Task-Management/config"

	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/gofiber/fiber/v2"
)

func InitNotificationRoutes(router *fiber.Router, app *app.Container, cfg config.Server) {

	notificationGroup := (*router).Group("/notifications", middlerwares.Auth([]byte(cfg.TokenSecret)))

	notificationGroup.Get("/", handlers.GetAllNotifications(app.NotificationService()))
	notificationGroup.Get("/unread", handlers.GetUnreadNotifications(app.NotificationService()))
	notificationGroup.Get("/:id", handlers.GetNotificationByID(app.NotificationService()))
	notificationGroup.Patch("/:id/read", handlers.ReadNotification(app.NotificationService()))
	notificationGroup.Patch("/:id/unread", handlers.UnReadNotification(app.NotificationService()))
	notificationGroup.Delete("/:id", handlers.DeleteNotification(app.NotificationService()))
}
