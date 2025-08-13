/**
 * Author : ngdangkietswe
 * Since  : 8/13/2025
 */

package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngdangkietswe/go-rabbitmq/internal/handlers"
	"log"
)

type NotificationRoutes struct {
	notificationHandler *handlers.NotificationHandler
}

func NewNotificationRoutes(notificationHandler *handlers.NotificationHandler) *NotificationRoutes {
	return &NotificationRoutes{
		notificationHandler: notificationHandler,
	}
}

func (r *NotificationRoutes) Register(router fiber.Router) {
	notification := router.Group("/notifications")

	notification.Post("/", r.notificationHandler.SendNotification)

	log.Printf("Notification routes registered successfully")
}
