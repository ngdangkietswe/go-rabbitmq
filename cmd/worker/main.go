/**
 * Author : ngdangkietswe
 * Since  : 8/12/2025
 */

package main

import (
	"encoding/json"
	"github.com/ngdangkietswe/go-rabbitmq/internal/config"
	"github.com/ngdangkietswe/go-rabbitmq/internal/models"
	"github.com/ngdangkietswe/go-rabbitmq/internal/services"
	"github.com/ngdangkietswe/go-rabbitmq/pkg/constants"
	"github.com/ngdangkietswe/go-rabbitmq/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.NewAppConfig("./configs")

	env := config.GetString("ENV", "local")

	appLogger := logger.NewAppLogger(env)

	rabbitMQUrl := config.GetString("RABBITMQ_URL", "amqp://admin:admin123@localhost:5672/")

	appLogger.Info("Connecting to RabbitMQ", zap.String("url", rabbitMQUrl))

	time.Sleep(10 * time.Second) // Simulate some startup delay

	rabbitMQ, err := services.NewRabbitMQService(rabbitMQUrl, appLogger)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer func(rabbitMQ *services.RabbitMQService) {
		if err := rabbitMQ.Close(); err != nil {
			// ignore
		}
	}(rabbitMQ)

	notificationService := services.NewNotificationService(appLogger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		if err := rabbitMQ.Close(); err != nil {
			// ignore
		}
		os.Exit(0)
	}()

	appLogger.Info("Starting RabbitMQ consumer...", zap.String("url", rabbitMQUrl))

	// consume messages from queue_notification
	if err := rabbitMQ.ConsumeMessages(string(constants.QueueNotification), func(delivery amqp.Delivery) error {
		var notification models.Notification

		if err := json.Unmarshal(delivery.Body, &notification); err != nil {
			if err := delivery.Nack(false, false); err != nil {
				appLogger.Error("Failed to nack message", zap.Error(err))
				return err
			}
		}

		return notificationService.ProcessNotification(&notification)
	}); err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	// consume messages from queue_log
	if err := rabbitMQ.ConsumeMessages(string(constants.QueueLog), func(delivery amqp.Delivery) error {
		appLogger.Info("Log message received", zap.ByteString("body", delivery.Body))
		return nil
	}); err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}
}
