/**
 * Author : ngdangkietswe
 * Since  : 8/12/2025
 */

package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ngdangkietswe/go-rabbitmq/internal/models"
	"github.com/ngdangkietswe/go-rabbitmq/internal/services"
)

type NotificationHandler struct {
	rabbitMQ *services.RabbitMQService
}

func NewNotificationHandler(rabbitMQ *services.RabbitMQService) *NotificationHandler {
	return &NotificationHandler{
		rabbitMQ: rabbitMQ,
	}
}

// SendNotification handles the request to send a notification.
// @Summary Send a notification
// @Description Send a notification to a recipient
// @Tags Notifications
// @Accept json
// @Produce json
// @Param notification body models.SendNotificationRequest true "Notification request"
// @Success 202 {object} models.SendNotificationResponse
// @Failure 400 {object} fiber.Error "Invalid request"
// @Failure 500 {object} fiber.Error "Internal server error"
// @Router /api/v1/notifications [post]
func (h *NotificationHandler) SendNotification(ctx *fiber.Ctx) error {
	var notificationRequest models.SendNotificationRequest

	if err := ctx.BodyParser(&notificationRequest); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if notificationRequest.Type == "" || notificationRequest.Recipient == "" || notificationRequest.Message == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing required fields")
	}

	if notificationRequest.Type != models.NotificationTypeEmail &&
		notificationRequest.Type != models.NotificationTypeSMS &&
		notificationRequest.Type != models.NotificationTypePush {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid notification type")
	}

	notification := &models.Notification{
		ID:        uuid.New().String(),
		Type:      notificationRequest.Type,
		Recipient: notificationRequest.Recipient,
		Title:     notificationRequest.Title,
		Message:   notificationRequest.Message,
		MetaData:  notificationRequest.MetaData,
		Status:    models.NotificationStatusPending,
	}

	if err := h.rabbitMQ.PublishMessage(notification); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to publish notification")
	}

	response := models.SendNotificationResponse{
		ID:     notification.ID,
		Status: notification.Status,
	}

	return ctx.Status(fiber.StatusAccepted).JSON(response)
}
