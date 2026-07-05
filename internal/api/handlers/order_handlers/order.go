package orderhandlers

import (
	"electra/internal/api/handlers/interfaces"
	"electra/pkg/logger"
)

type OrderHandler struct {
	orderService interfaces.OrderService
	logger       logger.Logger
}

func NewOrderHandler(orderService interfaces.OrderService, logger logger.Logger) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		logger:       logger,
	}
}
