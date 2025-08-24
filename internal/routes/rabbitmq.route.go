/**
 * Author : ngdangkietswe
 * Since  : 8/20/2025
 */

package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngdangkietswe/go-rabbitmq/internal/handlers"
)

type RabbitMQRoutes struct {
	rabbitMQHandler *handlers.RabbitMQHandler
}

func NewRabbitMQRoutes(rabbitMQHandler *handlers.RabbitMQHandler) *RabbitMQRoutes {
	return &RabbitMQRoutes{
		rabbitMQHandler: rabbitMQHandler,
	}
}

func (r *RabbitMQRoutes) Register(router fiber.Router) {
	rabbitMQ := router.Group("/rabbitmq")

	rabbitMQ.Get("/queues", r.rabbitMQHandler.GetListQueues)
	rabbitMQ.Get("/exchanges", r.rabbitMQHandler.GetListExchanges)
	rabbitMQ.Get("/bindings", r.rabbitMQHandler.GetListBindings)
}
