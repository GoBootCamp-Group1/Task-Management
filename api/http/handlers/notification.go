package handlers

import (
	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers/presenter"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/services"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/log"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrNotificationNotFound = fiber.NewError(fiber.StatusNotFound, "notification not found")
	ErrNoNotificationFound  = fiber.NewError(fiber.StatusNotFound, "no notification found")
)

// GetNotificationByID get a Notification
// @Summary Get Notification
// @Description gets a Notification
// @Tags Notification
// @Produce json
// @Param   id      path     string  true  "Notification ID"
// @Success 200 {object} domains.Notification
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /notifications/{id} [get]
// @Security ApiKeyAuth
func GetNotificationByID(NotificationService *services.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		userID, err := utils.GetUserID(c)
		if err != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		Notification, err := NotificationService.GetNotificationByID(c.Context(), id, userID)
		if err != nil {
			log.ErrorLog.Printf("Error getting Notification: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		if Notification == nil {
			log.ErrorLog.Printf("Error getting Notification: %v\n", ErrNotificationNotFound)
			return SendError(c, ErrNotificationNotFound, fiber.StatusNotFound)
		}
		log.InfoLog.Println("Notification loaded successfully")

		return SendSuccessResponse(
			c,
			"Successfully fetched.",
			presenter.NewNotificationPresenter(Notification),
		)
	}
}

// GetAllNotifications get all Notifications
// @Summary Get Notifications
// @Description gets Notifications for a user
// @Tags Notification
// @Produce json
// @Success 200 {array} presenter.NotificationPresenter
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /notifications [get]
// @Security ApiKeyAuth
func GetAllNotifications(NotificationService *services.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// init variables for pagination
		page, pageSize := PageAndPageSize(c)

		userID, err := utils.GetUserID(c)
		if err != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		Notifications, total, err := NotificationService.GetAllNotificationsList(c.Context(), userID, uint(page), uint(pageSize))
		if err != nil {
			log.ErrorLog.Printf("Error gettings Notifications: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		if len(Notifications) == 0 {
			log.ErrorLog.Printf("Error getting Notifications: %v\n", ErrNoNotificationFound)
			return SendError(c, ErrNoNotificationFound, fiber.StatusNotFound)
		}

		//generate response data
		NotificationPresenters := make([]*presenter.NotificationPresenter, len(Notifications))
		for i, Notification := range Notifications {
			NotificationPresenters[i] = presenter.NewNotificationPresenter(&Notification)
		}
		log.InfoLog.Println("Notifications loaded successfully")

		return SendSuccessPaginateResponse(
			c,
			"Successfully fetched.",
			NotificationPresenters,
			uint(page),
			uint(pageSize),
			total,
		)
	}
}

// GetUnreadNotifications get unread Notifications
// @Summary Get Unread Notifications
// @Description gets Unread Notifications for a user
// @Tags Notification
// @Produce json
// @Success 200 {array} presenter.NotificationPresenter
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /notifications/unread [get]
// @Security ApiKeyAuth
func GetUnreadNotifications(NotificationService *services.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// init variables for pagination
		page, pageSize := PageAndPageSize(c)

		userID, err := utils.GetUserID(c)
		if err != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		Notifications, total, err := NotificationService.GetUnReadNotificationsList(c.Context(), userID, uint(page), uint(pageSize))
		if err != nil {
			log.ErrorLog.Printf("Error gettings Notifications: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		if len(Notifications) == 0 {
			log.ErrorLog.Printf("Error getting Notifications: %v\n", ErrNoNotificationFound)
			return SendError(c, ErrNoNotificationFound, fiber.StatusNotFound)
		}

		//generate response data
		NotificationPresenters := make([]*presenter.NotificationPresenter, len(Notifications))
		for i, Notification := range Notifications {
			NotificationPresenters[i] = presenter.NewNotificationPresenter(&Notification)
		}
		log.InfoLog.Println("Notifications loaded successfully")

		return SendSuccessPaginateResponse(
			c,
			"Successfully fetched.",
			NotificationPresenters,
			uint(page),
			uint(pageSize),
			total,
		)
	}
}

// ReadNotification reads a Notification
// @Summary Read Notification
// @Description reads a Notification
// @Tags Notification
// @Accept  json
// @Produce json
// @Param   id      path     string  true  "Notification ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /notifications/{id}/read [patch]
// @Security ApiKeyAuth
func ReadNotification(NotificationService *services.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		userID, err := utils.GetUserID(c)
		if err != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		updatedNotification, err := NotificationService.ReadNotification(c.Context(), id, userID)
		if err != nil {
			log.ErrorLog.Printf("Error updating Notification: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}
		log.InfoLog.Println("Notification updated successfully")

		return SendSuccessResponse(
			c,
			"Successfully updated.",
			presenter.NewNotificationPresenter(updatedNotification),
		)
	}
}

// UnReadNotification unread a Notification
// @Summary UnRead Notification
// @Description unread a Notification
// @Tags Notification
// @Accept  json
// @Produce json
// @Param   id      path     string  true  "Notification ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /notifications/{id}/unread [patch]
// @Security ApiKeyAuth
func UnReadNotification(NotificationService *services.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		userID, err := utils.GetUserID(c)
		if err != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		updatedNotification, err := NotificationService.UnReadNotification(c.Context(), id, userID)
		if err != nil {
			log.ErrorLog.Printf("Error updating Notification: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}
		log.InfoLog.Println("Notification updated successfully")

		return SendSuccessResponse(
			c,
			"Successfully updated.",
			presenter.NewNotificationPresenter(updatedNotification),
		)
	}
}

// DeleteNotification delete a Notification
// @Summary Delete Notification
// @Description deleted a Notification
// @Tags Notification
// @Produce json
// @Param   id      path     string  true  "Notification ID"
// @Success 204
// @Failure 400
// @Failure 500
// @Router /notifications/{id} [delete]
// @Security ApiKeyAuth
func DeleteNotification(NotificationService *services.NotificationService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		userID, err := utils.GetUserID(c)
		if err != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		err = NotificationService.DeleteNotification(c.Context(), id, userID)
		if err != nil {
			log.ErrorLog.Printf("Error deleting Notification: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}
		msg := "Notification deleted successfully"
		log.InfoLog.Println(msg)

		return SendSuccessResponse(c, msg, id)
	}
}
