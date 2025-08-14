/**
 * Author : ngdangkietswe
 * Since  : 8/13/2025
 */

package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngdangkietswe/go-rabbitmq/internal/handlers"
	"time"
)

type AppRoutes struct {
	notificationHandler *handlers.NotificationHandler
}

func NewAppRoutes(notificationHandler *handlers.NotificationHandler) *AppRoutes {
	return &AppRoutes{
		notificationHandler: notificationHandler,
	}
}

func (r *AppRoutes) Register(app *fiber.App) {
	api := app.Group("/api/v1")

	// Heath check
	api.Get("/health", r.HealthCheck)

	// Notification routes
	notificationRoutes := NewNotificationRoutes(r.notificationHandler)
	notificationRoutes.Register(api)

	// Additional routes can be registered here
}

// HealthCheck handles the health check endpoint.
// @Summary Health check
// @Description Check if the notification service is running
// @Tags Health
// @Produce json
// @Success 200 {object} fiber.Map "Service is running"
// @Router /api/v1/health [get]
func (r *AppRoutes) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status":    "ok",
		"message":   "Notification service is running",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
