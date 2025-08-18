/**
 * Author : ngdangkietswe
 * Since  : 8/12/2025
 */

package main

import (
	"github.com/ngdangkietswe/go-rabbitmq/internal/config"
	"github.com/ngdangkietswe/go-rabbitmq/internal/models"
	"github.com/ngdangkietswe/go-rabbitmq/internal/services"
	"github.com/ngdangkietswe/go-rabbitmq/pkg/logger"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.NewAppConfig("./configs")

	env := config.GetString("ENV", "development")

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

	handler := func(notification *models.Notification) error {
		return notificationService.ProcessNotification(notification)
	}

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

	if err := rabbitMQ.ConsumeMessages(handler); err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}
}
