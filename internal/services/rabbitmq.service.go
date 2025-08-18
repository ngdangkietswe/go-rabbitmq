/**
 * Author : ngdangkietswe
 * Since  : 8/12/2025
 */

package services

import (
	"encoding/json"
	"fmt"
	"github.com/ngdangkietswe/go-rabbitmq/internal/models"
	"github.com/ngdangkietswe/go-rabbitmq/pkg/constants"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

type RabbitMQService struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	logger *zap.Logger
}

func NewRabbitMQService(url string, logger *zap.Logger) (*RabbitMQService, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	service := &RabbitMQService{conn: conn, ch: ch, logger: logger}

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
		string(constants.ExchangeNotification),
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
		string(constants.QueueNotification),
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := r.ch.QueueBind(
		string(constants.QueueNotification),
		string(constants.RoutingKeyNotification),
		string(constants.ExchangeNotification),
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
		string(constants.ExchangeNotification),
		string(constants.RoutingKeyNotification),
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

	r.logger.Info("Published notification", zap.String("id", notification.ID))

	return nil
}

func (r *RabbitMQService) ConsumeMessages(handler func(notification *models.Notification) error) error {
	msgs, err := r.ch.Consume(
		string(constants.QueueNotification),
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

	r.logger.Info("Consumer registered for notifications", zap.String("queue", string(constants.QueueNotification)))

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			var notification models.Notification
			if err := json.Unmarshal(msg.Body, &notification); err != nil {
				r.logger.Error("Failed to unmarshal notification", zap.Error(err))
				err := msg.Nack(false, false)
				if err != nil {
					return
				} // Reject the message
				continue
			}

			r.logger.Info("Received notification", zap.String("id", notification.ID))

			if err := handler(&notification); err != nil {
				r.logger.Error("Failed to process notification", zap.String("id", notification.ID), zap.Error(err))
				err := msg.Nack(false, false)
				if err != nil {
					return
				} // Reject the message
			} else {
				r.logger.Info("Processed notification successfully", zap.String("id", notification.ID))
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
