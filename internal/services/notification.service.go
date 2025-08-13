/**
 * Author : ngdangkietswe
 * Since  : 8/12/2025
 */

package services

import (
	"fmt"
	"github.com/ngdangkietswe/go-rabbitmq/internal/models"
	"log"
	"time"
)

type NotificationService struct {
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (ns *NotificationService) ProcessNotification(notification *models.Notification) error {
	notification.Status = models.NotificationStatusProcessing

	log.Printf("Processing %s notification: %s", notification.Type, notification.ID)

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

	log.Printf("Notification %s processed successfully", notification.ID)

	return nil
}

func (ns *NotificationService) sendEmail(notification *models.Notification) error {
	// Implement email sending logic here
	log.Printf("Sending email to %s: %s", notification.Recipient, notification.Message)
	return nil
}

func (ns *NotificationService) sendSMS(notification *models.Notification) error {
	// Implement SMS sending logic here
	log.Printf("Sending SMS to %s: %s", notification.Recipient, notification.Message)
	return nil
}

func (ns *NotificationService) sendPush(notification *models.Notification) error {
	// Implement push notification sending logic here
	log.Printf("Sending push notification to %s: %s", notification.Recipient, notification.Message)
	return nil
}
