/**
 * Author : ngdangkietswe
 * Since  : 8/12/2025
 */

package main

import (
	"github.com/ngdangkietswe/go-rabbitmq/internal/config"
	"github.com/ngdangkietswe/go-rabbitmq/internal/models"
	"github.com/ngdangkietswe/go-rabbitmq/internal/services"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	appConfig := config.NewAppConfig("./configs")

	rabbitMQUrl := appConfig.GetRabbitMQUrl()

	log.Println("Starting Notification Worker with RabbitMQ URL:", rabbitMQUrl)
	time.Sleep(10 * time.Second) // Simulate some startup delay

	rabbitMQ, err := services.NewRabbitMQService(rabbitMQUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer func(rabbitMQ *services.RabbitMQService) {
		if err := rabbitMQ.Close(); err != nil {
			// ignore
		}
	}(rabbitMQ)

	notificationService := services.NewNotificationService()

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

	log.Println("Waiting for messages...")

	if err := rabbitMQ.ConsumeMessages(handler); err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}
}
