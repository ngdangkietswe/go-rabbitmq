/**
 * Author : ngdangkietswe
 * Since  : 8/12/2025
 */

package services

import (
	"encoding/json"
	"fmt"
	"github.com/ngdangkietswe/go-rabbitmq/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

const (
	NotificationQueue      = "notification_queue"
	NotificationExchange   = "notification_exchange"
	NotificationRoutingKey = "notification"
)

type RabbitMQService struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQService(url string) (*RabbitMQService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	service := &RabbitMQService{conn: conn, ch: ch}

	if err := service.setupExchange(); err != nil {
		if err := service.Close(); err != nil {
			// ignore
		}
		return nil, err
	}

	if err := service.setupQueue(); err != nil {
		if err := service.Close(); err != nil {
			// ignore
		}
		return nil, err
	}

	return service, nil
}

func (r *RabbitMQService) setupExchange() error {
	if err := r.ch.ExchangeDeclare(
		NotificationExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	return nil
}

func (r *RabbitMQService) setupQueue() error {
	if _, err := r.ch.QueueDeclare(
		NotificationQueue,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := r.ch.QueueBind(
		NotificationQueue,
		NotificationRoutingKey,
		NotificationExchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return nil
}

func (r *RabbitMQService) PublishMessage(notification *models.Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	if err = r.ch.Publish(
		NotificationExchange,
		NotificationRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published notification: %s", notification.ID)

	return nil
}

func (r *RabbitMQService) ConsumeMessages(handler func(notification *models.Notification) error) error {
	msgs, err := r.ch.Consume(
		NotificationQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Waiting for messages in queue: %s", NotificationQueue)

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			var notification models.Notification
			if err := json.Unmarshal(msg.Body, &notification); err != nil {
				log.Printf("Failed to unmarshal message: %s", err)
				err := msg.Nack(false, false)
				if err != nil {
					return
				} // Reject the message
				continue
			}

			log.Printf("Received notification: %s", notification.ID)

			if err := handler(&notification); err != nil {
				log.Printf("Error processing notification: %s", err)
				err := msg.Nack(false, false)
				if err != nil {
					return
				} // Reject the message
			} else {
				log.Printf("Successfully processed notification: %s", notification.ID)
				err := msg.Ack(false)
				if err != nil {
					return
				} // Acknowledge the message
			}
		}
	}()

	<-forever

	return nil
}

func (r *RabbitMQService) Close() error {
	if err := r.ch.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}
	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}
