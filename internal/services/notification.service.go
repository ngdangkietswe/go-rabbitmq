/**
 * Author : ngdangkietswe
 * Since  : 8/12/2025
 */

package services

import (
	"fmt"
	"github.com/ngdangkietswe/go-rabbitmq/internal/models"
	"go.uber.org/zap"
	"time"
)

type NotificationService struct {
	logger *zap.Logger
}

func NewNotificationService(logger *zap.Logger) *NotificationService {
	return &NotificationService{
		logger: logger,
	}
}

func (ns *NotificationService) ProcessNotification(notification *models.Notification) error {
	notification.Status = models.NotificationStatusProcessing

	ns.logger.Info("Processing notification", zap.String("id", notification.ID), zap.String("type", string(notification.Type)), zap.String("recipient", notification.Recipient))

	switch notification.Type {
	case models.NotificationTypeEmail:
		if err := ns.sendEmail(notification); err != nil {
			notification.Status = models.NotificationStatusFailed
			notification.Error = err.Error()
			return fmt.Errorf("failed to send email notification %s: %w", notification.ID, err)
		}
	case models.NotificationTypeSMS:
		if err := ns.sendSMS(notification); err != nil {
			notification.Status = models.NotificationStatusFailed
			notification.Error = err.Error()
			return fmt.Errorf("failed to send SMS notification %s: %w", notification.ID, err)
		}
	case models.NotificationTypePush:
		if err := ns.sendPush(notification); err != nil {
			notification.Status = models.NotificationStatusFailed
			notification.Error = err.Error()
			return fmt.Errorf("failed to send push notification %s: %w", notification.ID, err)
		}
	}

	now := time.Now()

	notification.Status = models.NotificationStatusSent
	notification.SentAt = &now

	ns.logger.Info("Notification sent successfully", zap.String("id", notification.ID), zap.String("type", string(notification.Type)), zap.String("recipient", notification.Recipient))

	return nil
}

func (ns *NotificationService) sendEmail(notification *models.Notification) error {
	// Implement email sending logic here
	ns.logger.Info("Sending email", zap.String("recipient", notification.Recipient), zap.String("message", notification.Message))
	return nil
}

func (ns *NotificationService) sendSMS(notification *models.Notification) error {
	// Implement SMS sending logic here
	ns.logger.Info("Sending SMS", zap.String("recipient", notification.Recipient), zap.String("message", notification.Message))
	return nil
}

func (ns *NotificationService) sendPush(notification *models.Notification) error {
	// Implement push notification sending logic here
	ns.logger.Info("Sending push", zap.String("recipient", notification.Recipient), zap.String("message", notification.Message))
	return nil
}
