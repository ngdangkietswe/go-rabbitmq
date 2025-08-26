/**
 * Author : ngdangkietswe
 * Since  : 8/12/2025
 */

package services

import (
	"encoding/json"
	"fmt"
	"github.com/ngdangkietswe/go-rabbitmq/internal/config"
	"github.com/ngdangkietswe/go-rabbitmq/internal/models"
	"github.com/ngdangkietswe/go-rabbitmq/pkg/constants"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type RabbitMQService struct {
	conn                   *amqp.Connection
	ch                     *amqp.Channel
	logger                 *zap.Logger
	rabbitMQManagementAddr string
	rabbitMQVHost          string
	rabbitMQUsername       string
	rabbitMQPassword       string
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

	service := &RabbitMQService{
		conn:                   conn,
		ch:                     ch,
		logger:                 logger,
		rabbitMQManagementAddr: config.GetString("RABBITMQ_MANAGEMENT_ADDRESS", "http://localhost:15672"),
		rabbitMQVHost:          config.GetString("RABBITMQ_VHOST", "%2F"),
		rabbitMQUsername:       config.GetString("RABBITMQ_USERNAME", "admin"),
		rabbitMQPassword:       config.GetString("RABBITMQ_PASSWORD", "admin123"),
	}

	exchanges := []*models.ExchangeConfig{
		{
			Name:       string(constants.ExchangeNotification),
			Type:       "direct",
			Durable:    true,
			AutoDelete: false,
		},
		{
			Name:       string(constants.ExchangeLog),
			Type:       "fanout",
			Durable:    true,
			AutoDelete: false,
		},
	}

	if err := service.setupExchanges(exchanges); err != nil {
		logger.Error("Failed to setup exchanges", zap.Error(err))
		return nil, err
	}

	queues := []*models.QueueConfig{
		{
			Name:       string(constants.QueueNotification),
			Durable:    true,
			AutoDelete: false,
		},
		{
			Name:       string(constants.QueueLog),
			Durable:    true,
			AutoDelete: false,
		},
	}

	if err := service.setupQueues(queues); err != nil {
		logger.Error("Failed to setup queues", zap.Error(err))
		return nil, err
	}

	bindings := []*models.BindingConfig{
		{
			Exchange:   string(constants.ExchangeNotification),
			Queue:      string(constants.QueueNotification),
			RoutingKey: string(constants.RoutingKeyNotification),
		},
		{
			Exchange:   string(constants.ExchangeLog),
			Queue:      string(constants.QueueLog),
			RoutingKey: "",
		},
	}

	if err := service.setupBindings(bindings); err != nil {
		logger.Error("Failed to setup bindings", zap.Error(err))
		return nil, err
	}

	return service, nil
}

func (r *RabbitMQService) setupExchanges(exchanges []*models.ExchangeConfig) error {
	lo.ForEach(exchanges, func(exchange *models.ExchangeConfig, _ int) {
		if err := r.ch.ExchangeDeclare(
			exchange.Name,
			exchange.Type,
			exchange.Durable,
			exchange.AutoDelete,
			false,
			false,
			nil,
		); err != nil {
			r.logger.Error("Failed to declare exchange", zap.String("exchange", exchange.Name), zap.Error(err))
		} else {
			r.logger.Info("Declared exchange", zap.String("exchange", exchange.Name))
		}
	})

	return nil
}

func (r *RabbitMQService) setupQueues(queues []*models.QueueConfig) error {
	lo.ForEach(queues, func(queue *models.QueueConfig, _ int) {
		if _, err := r.ch.QueueDeclare(
			queue.Name,
			queue.Durable,
			queue.AutoDelete,
			false,
			false,
			nil,
		); err != nil {
			r.logger.Error("Failed to declare queue", zap.String("queue", queue.Name), zap.Error(err))
		} else {
			r.logger.Info("Declared queue", zap.String("queue", queue.Name))
		}
	})

	return nil
}

func (r *RabbitMQService) setupBindings(bindings []*models.BindingConfig) error {
	lo.ForEach(bindings, func(binding *models.BindingConfig, _ int) {
		if err := r.ch.QueueBind(
			binding.Queue,
			binding.RoutingKey,
			binding.Exchange,
			false,
			nil,
		); err != nil {
			r.logger.Error("Failed to bind queue", zap.String("queue", binding.Queue), zap.String("exchange", binding.Exchange), zap.Error(err))
		} else {
			r.logger.Info("Bound queue", zap.String("queue", binding.Queue), zap.String("exchange", binding.Exchange))
		}
	})

	return nil
}

func (r *RabbitMQService) PublishMessage(exchange, routingKey string, notification *models.Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	if err = r.ch.Publish(
		exchange,
		routingKey,
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

func (r *RabbitMQService) ConsumeMessages(queue string, handler func(delivery amqp.Delivery) error) error {
	msgs, err := r.ch.Consume(
		queue, "", false, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("consume queue %s: %w", queue, err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			if err := handler(msg); err != nil {
				r.logger.Error("Failed to process message", zap.Error(err))
			}
		}
	}()

	r.logger.Info("Consumer registered", zap.String("queue", queue))

	<-forever

	return nil
}

func (r *RabbitMQService) GetListQueues() ([]map[string]any, error) {
	url := fmt.Sprintf("%s/api/queues/%s", r.rabbitMQManagementAddr, r.rabbitMQVHost)

	return r.getAndDecodeAPI(url)
}

func (r *RabbitMQService) GetListExchanges() ([]map[string]any, error) {
	url := fmt.Sprintf("%s/api/exchanges/%s", r.rabbitMQManagementAddr, r.rabbitMQVHost)

	return r.getAndDecodeAPI(url)
}

func (r *RabbitMQService) GetListBindings() ([]map[string]any, error) {
	url := fmt.Sprintf("%s/api/bindings/%s", r.rabbitMQManagementAddr, r.rabbitMQVHost)

	return r.getAndDecodeAPI(url)
}

func (r *RabbitMQService) getAndDecodeAPI(url string) ([]map[string]any, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.SetBasicAuth(r.rabbitMQUsername, r.rabbitMQPassword)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			r.logger.Error("failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	var data []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode response body: %w", err)
	}

	return data, nil
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
