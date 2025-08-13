/**
 * Author : ngdangkietswe
 * Since  : 8/12/2025
 */

package models

import "time"

type NotificationType string

const (
	NotificationTypeEmail NotificationType = "email"
	NotificationTypeSMS   NotificationType = "sms"
	NotificationTypePush  NotificationType = "push"
)

type NotificationStatus string

const (
	NotificationStatusPending    NotificationStatus = "pending"
	NotificationStatusProcessing NotificationStatus = "processing"
	NotificationStatusSent       NotificationStatus = "sent"
	NotificationStatusFailed     NotificationStatus = "failed"
)

type Notification struct {
	ID        string                 `json:"id"`
	Type      NotificationType       `json:"type"`
	Recipient string                 `json:"recipient"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	MetaData  map[string]interface{} `json:"meta_data,omitempty"` // Optional metadata for additional information
	Status    NotificationStatus     `json:"status"`
	SentAt    *time.Time             `json:"sent_at,omitempty"`
	Error     string                 `json:"error,omitempty"` // Optional error message if the notification fails
}

type SendNotificationRequest struct {
	Type      NotificationType       `json:"type"`
	Recipient string                 `json:"recipient"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	MetaData  map[string]interface{} `json:"meta_data,omitempty"` // Optional metadata for additional information
}

type SendNotificationResponse struct {
	ID     string             `json:"id"`
	Status NotificationStatus `json:"status"`
}
