/**
 * Author : ngdangkietswe
 * Since  : 8/20/2025
 */

package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngdangkietswe/go-rabbitmq/internal/services"
	"go.uber.org/zap"
)

type RabbitMQHandler struct {
	rabbitMQ *services.RabbitMQService
	logger   *zap.Logger
}

func NewRabbitMQHandler(rabbitMQ *services.RabbitMQService, logger *zap.Logger) *RabbitMQHandler {
	return &RabbitMQHandler{
		rabbitMQ: rabbitMQ,
		logger:   logger,
	}
}

// GetListQueues godoc
// @Summary Get list of RabbitMQ queues
// @Description Retrieve a list of all RabbitMQ queues
// @Tags RabbitMQ
// @Accept json
// @Produce json
// @Success 200 {array} interface{} "List of RabbitMQ queues"
// @Failure 500 {object} fiber.Error "Internal server error"
// @Router /api/v1/rabbitmq/queues [get]
func (h *RabbitMQHandler) GetListQueues(ctx *fiber.Ctx) error {
	queues, err := h.rabbitMQ.GetListQueues()
	if err != nil {
		h.logger.Error("Failed to get queues", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get queues")
	}

	return ctx.Status(fiber.StatusOK).JSON(queues)
}

// GetListExchanges godoc
// @Summary Get list of RabbitMQ exchanges
// @Description Retrieve a list of all RabbitMQ exchanges
// @Tags RabbitMQ
// @Accept json
// @Produce json
// @Success 200 {array} interface{} "List of RabbitMQ exchanges"
// @Failure 500 {object} fiber.Error "Internal server error"
// @Router /api/v1/rabbitmq/exchanges [get]
func (h *RabbitMQHandler) GetListExchanges(ctx *fiber.Ctx) error {
	exchanges, err := h.rabbitMQ.GetListExchanges()
	if err != nil {
		h.logger.Error("Failed to get exchanges", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get exchanges")
	}

	return ctx.Status(fiber.StatusOK).JSON(exchanges)
}

// GetListBindings godoc
// @Summary Get list of RabbitMQ bindings
// @Description Retrieve a list of all RabbitMQ bindings
// @Tags RabbitMQ
// @Accept json
// @Produce json
// @Success 200 {array} interface{} "List of RabbitMQ bindings"
// @Failure 500 {object} fiber.Error "Internal server error"
// @Router /api/v1/rabbitmq/bindings [get]
func (h *RabbitMQHandler) GetListBindings(ctx *fiber.Ctx) error {
	bindings, err := h.rabbitMQ.GetListBindings()
	if err != nil {
		h.logger.Error("Failed to get bindings", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get bindings")
	}

	return ctx.Status(fiber.StatusOK).JSON(bindings)
}
