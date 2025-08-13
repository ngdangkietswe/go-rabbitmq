/**
 * Author : ngdangkietswe
 * Since  : 8/12/2025
 */

package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ngdangkietswe/go-rabbitmq/internal/config"
	"github.com/ngdangkietswe/go-rabbitmq/internal/handlers"
	"github.com/ngdangkietswe/go-rabbitmq/internal/services"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.NewAppConfig("./configs")

	httpPort := config.GetInt("HTTP_PORT", 3000)
	rabbitMQUrl := config.GetString("RABBITMQ_URL", "amqp://admin:admin123@localhost:5672/")

	log.Println("Starting Notification API with RabbitMQ URL:", rabbitMQUrl)
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

	app := fiber.New(fiber.Config{
		AppName: "Notification Service v1.0",
	})

	app.Use(logger.New())
	app.Use(cors.New())

	api := app.Group("/api/v1")

	notificationHandler := handlers.NewNotificationHandler(rabbitMQ)

	api.Get("/health", notificationHandler.HealthCheck)
	api.Post("/notifications", notificationHandler.SendNotification)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Error shutting down: %v", err)
		}
	}()

	log.Printf("Starting server on port %d...", httpPort)

	if err := app.Listen(fmt.Sprintf(":%d", httpPort)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
