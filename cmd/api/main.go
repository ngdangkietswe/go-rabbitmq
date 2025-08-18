/**
 * Author : ngdangkietswe
 * Since  : 8/12/2025
 */

package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/ngdangkietswe/go-rabbitmq/docs" // swagger docs
	"github.com/ngdangkietswe/go-rabbitmq/internal/config"
	"github.com/ngdangkietswe/go-rabbitmq/internal/handlers"
	"github.com/ngdangkietswe/go-rabbitmq/internal/middlewares"
	"github.com/ngdangkietswe/go-rabbitmq/internal/routes"
	"github.com/ngdangkietswe/go-rabbitmq/internal/services"
	"github.com/ngdangkietswe/go-rabbitmq/pkg/logger"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Notification Service API
// @version 1.0
// @description This is a notification service API with RabbitMQ integration
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email ngdangkietswe@yopmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:3000
// @BasePath /api/v1
// @schemes http https
func main() {
	config.NewAppConfig("./configs")

	env := config.GetString("ENV", "development")
	httpPort := config.GetInt("HTTP_PORT", 3000)
	rabbitMQUrl := config.GetString("RABBITMQ_URL", "amqp://admin:admin123@localhost:5672/")

	appLogger := logger.NewAppLogger(env)

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

	app := fiber.New(fiber.Config{
		AppName: "Notification Service v1.0",
	})

	app.Use(middlewares.NewLogger())
	app.Use(middlewares.NewCORS())

	app.Get("/swagger/*", swagger.HandlerDefault)

	notificationHandler := handlers.NewNotificationHandler(rabbitMQ, appLogger)

	appRoutes := routes.NewAppRoutes(notificationHandler)
	appRoutes.Register(app)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Error shutting down: %v", err)
		}
	}()

	appLogger.Info("Starting server", zap.Int("port", httpPort))

	if err := app.Listen(fmt.Sprintf(":%d", httpPort)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
